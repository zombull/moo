package server

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/labstack/echo"

	"github.com/zombull/moo/bug"
	"github.com/zombull/moo/database"
	"github.com/zombull/moo/moonboard"
)

type KeyValueStore struct {
	cache  string
	server string
	data   map[string][]byte
	sums   map[string][]byte
}

func NewStore(cache, server string) *KeyValueStore {
	cache = path.Join(cache, "moonboard")
	server = path.Join(server, "moonboard")

	s := KeyValueStore{
		cache:  cache,
		server: server,
		data:   make(map[string][]byte),
		sums:   make(map[string][]byte),
	}

	infos, err := ioutil.ReadDir(s.cache)
	bug.OnError(err)

	for _, fi := range infos {
		if fi.Mode().IsRegular() {
			name := path.Join(s.cache, fi.Name())

			if strings.HasSuffix(fi.Name(), ".json") {
				s.data[strings.TrimSuffix(fi.Name(), ".json")], err = ioutil.ReadFile(name)
				bug.OnError(err)
			} else if strings.HasSuffix(fi.Name(), ".md5") {
				s.sums[strings.TrimSuffix(fi.Name(), ".md5")], err = ioutil.ReadFile(name)
				bug.OnError(err)
			}
		}
	}

	return &s
}

const internalServerError = "I'm freakin' out, man!  Please try again at a later time."

func (s *KeyValueStore) get(c echo.Context, key, notFound string) error {
	val, ok := s.data[key]
	if !ok && len(notFound) > 0 {
		return echo.NewHTTPError(http.StatusNotFound, notFound)
	} else if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, internalServerError)
	}
	return c.JSONBlob(http.StatusOK, val)
}

func (s *KeyValueStore) getInternal(key string) func(c echo.Context) error {
	return func(c echo.Context) error {
		return s.get(c, key, "")
	}
}

func (s *KeyValueStore) getValue(host string) func(c echo.Context) error {
	return func(c echo.Context) error {
		key := c.Param("key")
		return s.get(c, host+"."+key, fmt.Sprintf("Did not find any %s for '%s'", key, host))
	}
}

func sanitize(s string) string {
	return strings.ToLower(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s))
}

func checksum(b []byte) []byte {
	return []byte(fmt.Sprintf("%x", md5.Sum(b)))
}

func (s *KeyValueStore) export(p, k, d string, v interface{}, sums map[string][]byte) {
	b, err := json.Marshal(v)
	bug.OnError(err)

	key := p + "." + k
	s.data[key] = b
	err = ioutil.WriteFile(path.Join(d, key + ".json"), b, 0644)
	bug.OnError(err)

	sums[k] = checksum(b)
	s.sums[key] = sums[k]
	err = ioutil.WriteFile(path.Join(d, key + ".md5"), sums[k], 0644)
	bug.OnError(err)
}

type holds struct {
	Holds  string `json:"h"`
	Start  int    `json:"s"`
	Finish int    `json:"f"`
}

// moonEntry holds the actual data for a problem or a setter.
// Keep them in the same object even though many of the fields
// are unique to one or the other, as some parts of the code
// works with both problems and setters and so the JSON names
// need to be identical for common fields and distinct for unique
// fields.
type moonEntry struct {
	Url           string    `json:"u"`
	Name          string    `json:"n"`
	LowerCaseName string    `json:"l"`
	Id            int       `json:"i"` // index into [Problem|Setter]Data, not database ID or moonboard ID
	Date          string    `json:"d,omitempty"`
	Nickname      string    `json:"k,omitempty"`
	Holds         string    `json:"h,omitempty"`
	Problems      []int     `json:"p,omitempty"`
	Setter        int       `json:"r,omitempty"`
	Grade         string    `json:"g,omitempty"`
	Stars         uint      `json:"s,omitempty"`
	Ascents       uint      `json:"a,omitempty"`
	Benchmark     bool      `json:"b,omitempty"`
	Comment       string    `json:"c,omitempty"`
	MoonId        uint      `json:"-"`
	RawDate       time.Time `json:"-"`
}

type moonIndex struct {
	Problems []moonEntry
	Setters  []moonEntry
}
type moonData struct {
	Index    moonIndex
	Problems map[string]int
	Setters  map[string]int
	Images   []string
}

func getProblemUrl(s string) string {
	ss := strings.Split(strings.Trim(s, "/"), "/")
	s = ss[len(ss)-1]
	bug.On(len(s) == 0, fmt.Sprintf("%d %v", len(ss), ss))
	bug.On(s != strings.ToLower(s), fmt.Sprintf("Moonboard has a case sensitive URL? '%s' != '%s", s, strings.ToLower(s)))
	return s
}

func getSetterUrl(s string) string {
	return "s/" + url.PathEscape(sanitize(s))
}

func (s *KeyValueStore) Update(d *database.Database, area *database.Area) {
	cragId := moonboard.CragId(d)

	setters := d.GetSetters(cragId)
	bug.On(len(setters) == 0, fmt.Sprintf("No moonboard setters found: %d", cragId))

	routes := d.GetRoutes(cragId, area.Id)
	if (len(routes) == 0 && area.Name == "MoonBoard 2019") {
		return
	}

	bug.On(len(routes) == 0, fmt.Sprintf("No moonboard routes found: %d %d", cragId, area.Id))

	md := moonData{
		Index: moonIndex{
			Problems: make([]moonEntry, len(routes)),
			Setters:  make([]moonEntry, 0, len(setters)),
		},
		Problems: make(map[string]int),
		Setters:  make(map[string]int),
		Images:   make([]string, 281),
	}

	for _, r := range setters {
		e := moonEntry{
			Url:           getSetterUrl(r.Name),
			Name:          r.Name,
			Nickname:      r.Nickname,
			LowerCaseName: strings.ToLower(r.Name),
			Problems:      make([]int, 0),
		}

		// Like usual, Moonboard doesn't sanitize their data and so there
		// are "duplicate" setters that are actually the same person, just
		// with different capitalization of their name.  Assume all such
		// collisions are cases where it's a single setter and only insert
		// a new setter if they have a unique URL.  This is why append is
		// used instead of directly indexing, and is also why Setters is
		// created with a length of 0.
		if _, ok := md.Setters[e.Url]; !ok {
			e.Id = len(md.Index.Setters)
			md.Setters[e.Url] = e.Id
			md.Index.Setters = append(md.Index.Setters, e)
		}
	}

	sort.Slice(routes, func(i, j int) bool {
		p1 := routes[i]
		p2 := routes[j]

		// Note that the return is inverted from what might be expected
		// by a "Less" function, as we effectively want a reverse sort,
		// e.g. higher stars and ascents at the front of the list.  And
		// don't forget that Pitches is actualy Ascents, we're sorting
		// routes from the database, not the Moonboard specific problems.
		return (p1.Stars * p1.Stars * p1.Pitches) > (p2.Stars * p2.Stars * p2.Pitches)
	})

	for i, r := range routes {
		sn := d.GetSetter(r.SetterId).Name
		setter, ok := md.Setters[getSetterUrl(d.GetSetter(r.SetterId).Name)]
		bug.On(!ok, fmt.Sprintf("Moonboard problem has undefined setter: %s", sn))

		e := moonEntry{
			Url:           getProblemUrl(r.Url),
			Name:          r.Name,
			LowerCaseName: strings.ToLower(r.Name),
			Date:          r.Date.Format("2006-01-02"), // 'yyyy-MM-dd'
			Setter:        setter,
			Grade:         r.Grade,
			Stars:         r.Stars,
			Id:            i,
			Ascents:       r.Pitches,
			Benchmark:     r.Benchmark,
			Comment:       r.Comment,
			MoonId:        r.Length,
			RawDate:       r.Date,
		}

		holdMap := make(map[string]bool)
		h2 := d.GetHolds(r.Id)

		start := make([]string, 0)
		finish := make([]string, 0)
		intermediate := make([]string, 0)
		for _, v := range h2.Holds {
			h := string(v[1:])
			_, ok := holdMap[h]
			bug.On(ok, fmt.Sprintf("Duplicate hold %s in moonboard problem %s", h, e.Name))
			holdMap[h] = true
			if string(v[0]) == "s" {
				start = append(start, h)
			} else if string(v[0]) == "f" {
				finish = append(finish, h)
			} else {
				intermediate = append(intermediate, h)
			}
		}

		e.Holds = strings.Join(start, "") + "," + strings.Join(intermediate, "") + "," + strings.Join(finish, "")

		bug.On(len(start) == 0, fmt.Sprintf("%s: No start hold found", r.Name))
		bug.On(len(finish) == 0, fmt.Sprintf("%s: No finish hold found", r.Name))

		if existing, ok := md.Problems[e.Url]; ok {
			ex := &md.Index.Problems[existing]
			before := ex.RawDate.Before(e.RawDate)
			if ex.RawDate.Equal(e.RawDate) {
				bug.On(ex.MoonId == e.MoonId, fmt.Sprintf("Duplicate Moonboard problem: %s, %d and %d\n", e.Url, ex.MoonId, e.MoonId))
				before = ex.MoonId < e.MoonId
			}
			if before {
				e.Url = fmt.Sprintf("%d-%s", e.MoonId, e.Url)
				_, ok = md.Problems[e.Url]
				bug.On(ok, fmt.Sprintf("Duplicate Moonboard problem URL: %s", e.Url))
				fmt.Printf("Duplicate Moonboard problem, new URL: %s\n", e.Url)
			} else {
				ex.Url = fmt.Sprintf("%d-%s", ex.MoonId, ex.Url)
				_, ok = md.Problems[ex.Url]
				bug.On(ok, fmt.Sprintf("Duplicate Moonboard problem URL: %s", ex.Url))
				md.Problems[ex.Url] = existing
				fmt.Printf("Duplicate Moonboard problem, updated existing URL: %s\n", ex.Url)
			}
		}
		md.Problems[e.Url] = i

		md.Index.Problems[i] = e
		md.Index.Setters[setter].Problems = append(md.Index.Setters[setter].Problems, i)
	}

	imgDir := path.Join(s.server, "img")

	// Plastic holds are 1 - 40 and 51 - 201, but hold 201 is shifted to
	// 41 so that the wood holds can sanely start at 201.  The board is
	// slotted in at 0, thus 0..280 images, with holes at 42..50.
	for i := 0; i < 281; i++ {
		if i > 41 && i < 50 {
			continue
		}
		n := "board"
		if (i > 200) {
			n = "w" + strconv.Itoa(i - 200)
		} else if i > 0 {
			n = strconv.Itoa(i)
		}

		img, err := ioutil.ReadFile(path.Join(imgDir, n+".png"))
		bug.OnError(err)

		md.Images[i] = base64.StdEncoding.EncodeToString(img)
	}

	sums := make(map[string][]byte)

	prefix := strings.Replace(strings.ToLower(area.Name), " ", "", -1)
	s.export(prefix, "index.problems", s.cache, md.Index.Problems, sums)
	s.export(prefix, "index.setters", s.cache, md.Index.Setters, sums)
	s.export(prefix, "images", s.cache, md.Images, sums)
	s.export(prefix, "problems", s.cache, md.Problems, sums)
	s.export(prefix, "setters", s.cache, md.Setters, sums)

	sumj, err := json.Marshal(sums)
	bug.OnError(err)

	err = ioutil.WriteFile(path.Join(s.server, prefix + "md5.json"), sumj, 0644)
	bug.OnError(err)
}
