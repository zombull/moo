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
	data   map[string][]byte
	sums   map[string][]byte
}

func NewStore(cache string) *KeyValueStore {
	s := KeyValueStore{
		cache:  cache,
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

func (s *KeyValueStore) export(k, d string, v interface{}) {
	b, err := json.Marshal(v)
	bug.OnError(err)
	s.data[k] = b
	s.sums[k] = checksum(b)
	err = ioutil.WriteFile(path.Join(d, k+".json"), b, 0644)
	bug.OnError(err)
	err = ioutil.WriteFile(path.Join(d, k+".md5"), checksum(b), 0644)
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
	bug.On(s != strings.ToLower(s), "Moonboard has a case sensitive URL?")
	return s
}

func getSetterUrl(s string) string {
	return "s/" + url.PathEscape(sanitize(s))
}

func (s *KeyValueStore) Update(d *database.Database, server string) {
	setters := d.GetSetters(moonboard.Id(d))
	bug.On(len(setters) == 0, fmt.Sprintf("No moonboard setters found: %d", moonboard.Id(d)))

	routes := d.GetRoutes(moonboard.Id(d), moonboard.SetId(d))
	bug.On(len(routes) == 0, fmt.Sprintf("No moonboard routes found: %d", moonboard.Id(d)))

	md := moonData{
		Index: moonIndex{
			Problems: make([]moonEntry, len(routes)),
			Setters:  make([]moonEntry, 0, len(setters)),
		},
		Problems: make(map[string]int),
		Setters:  make(map[string]int),
		Images:   make([]string, 233),
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

	imgDir := path.Join(server, "img")
	for i := 0; i < 233; i++ {
		if i > 41 && i < 50 {
			continue
		}
		n := "board"
		if i > 0 {
			n = strconv.Itoa(i)
		}
		img, err := ioutil.ReadFile(path.Join(imgDir, n+".png"))
		bug.OnError(err)

		md.Images[i] = base64.StdEncoding.EncodeToString(img)
	}

	s.export("moonboard.index.problems", s.cache, md.Index.Problems)
	s.export("moonboard.index.setters", s.cache, md.Index.Setters)
	s.export("moonboard.images", s.cache, md.Images)
	s.export("moonboard.problems", s.cache, md.Problems)
	s.export("moonboard.setters", s.cache, md.Setters)
}
