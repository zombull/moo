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
		if unicode.IsSpace(r) || r == '\'' || r == '"' {
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

// moonProblem and moonSetter hold the actual data for a problem or a
// setter.  Define a common object purely to detect overlap, which must
// be avoided as some parts of the code work with both problems and
// setters.  The JSON names need to be identical for common fields and
// distinct for unique fields.  Using omitempty gets dangerous because
// valid, required values, e.g. '0', will get dropped :-(
type moonEntryCommon struct {
	Url           string    `json:"u"`
	Name          string    `json:"n"`
	LowerCaseName string    `json:"l"`
	Id            int       `json:"i"` // index into [Problem|Setter]Data, not database ID or moonboard ID

	Nickname      string    `json:"k,omitempty"`
	Problems      []int     `json:"p"`
	Date          string    `json:"d"`
	Holds         string    `json:"h"`
	Setter        int       `json:"r"`
	Grade         string    `json:"g"`
	Stars         uint      `json:"s"`
	Ascents       uint      `json:"a"`
	Benchmark     bool      `json:"b,omitempty"`
	Comment       string    `json:"c,omitempty"`
	MoonId        uint      `json:"-"`
	RawDate       time.Time `json:"-"`
}

type moonSetter struct {
	Url           string    `json:"u"`
	Name          string    `json:"n"`
	LowerCaseName string    `json:"l"`
	Id            int       `json:"i"` // index into [Problem|Setter]Data, not database ID or moonboard ID

	Nickname      string    `json:"k,omitempty"`
	Problems      []int     `json:"p"`
}

type moonProblem struct {
	Url           string    `json:"u"`
	Name          string    `json:"n"`
	LowerCaseName string    `json:"l"`
	Id            int       `json:"i"` // index into [Problem|Setter]Data, not database ID or moonboard ID

	Date          string    `json:"d"`
	Holds         string    `json:"h"`
	Setter        int       `json:"r"`
	Grade         string    `json:"g"`
	Stars         uint      `json:"s"`
	Ascents       uint      `json:"a"`
	Benchmark     bool      `json:"b,omitempty"`
	Comment       string    `json:"c,omitempty"`
	MoonId        uint      `json:"-"`
	RawDate       time.Time `json:"-"`
}

type moonIndex struct {
	Problems []moonProblem
	Setters  []moonSetter
}
type moonData struct {
	Index    moonIndex
	Problems map[string]int
	Setters  map[string]int
	Images   []string
}

func getSetterUrl(s string) string {
	return "s/" + url.PathEscape(sanitize(s))
}

func (s *KeyValueStore) Update(d *database.Database, set *database.Set) {
	setters := d.GetSetters()
	bug.On(len(setters) == 0, fmt.Sprintf("No moonboard setters found"))

	problems := d.GetProblems(set.Id)
	bug.On(len(problems) == 0, fmt.Sprintf("No moonboard problems found for '%s'", set.Name))

	md := moonData{
		Index: moonIndex{
			Problems: make([]moonProblem, len(problems)),
			Setters:  make([]moonSetter, 0, len(setters)),
		},
		Problems: make(map[string]int),
		Setters:  make(map[string]int),
		Images:   make([]string, 281),
	}

	for _, r := range setters {
		e := moonSetter{
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
		// created with a length of 0.  Note, this results in the setter
		// indices not being stable across updates due to the database
		// query not returning results sorted by setter Id (though the
		// results are stable for a given instance).
		if _, ok := md.Setters[e.Url]; !ok {
			e.Id = len(md.Index.Setters)
			md.Setters[e.Url] = e.Id
			md.Index.Setters = append(md.Index.Setters, e)
		}
	}

	for _, r := range problems {
		// Deduct two stars (give or take) in order to combat inflation.
		// Do this before sorting, which takes the number of stars into
		// account
		if r.Stars < 2 {
			r.Stars = 0;
		} else if r.Stars < 4 {
			r.Stars = 1;
		} else {
			r.Stars = r.Stars - 2;
		}
	}

	sort.Slice(problems, func(i, j int) bool {
		p1 := problems[i]
		p2 := problems[j]

		// Note that the return is inverted from what might be expected
		// by a "Less" function, as we effectively want a reverse sort,
		// e.g. higher stars and ascents at the front of the list.  And
		// we're sorting problems from the database, not the Moonboard
		// specific problems.
		return (p1.Stars * p1.Stars * p1.Ascents) > (p2.Stars * p2.Stars * p2.Ascents)
	})

	for i, r := range problems {
		sn := d.GetSetter(r.SetterId).Name
		setter, ok := md.Setters[getSetterUrl(d.GetSetter(r.SetterId).Name)]
		bug.On(!ok, fmt.Sprintf("Moonboard problem has undefined setter: %s", sn))

		var date time.Time = time.Unix(r.Date, 0)

		e := moonProblem{
			Url:           strconv.Itoa(int(r.MoonId)),
			Name:          r.Name,
			LowerCaseName: strings.ToLower(r.Name),
			Date:          date.Format("2006-01-02"), // 'yyyy-MM-dd'
			Setter:        setter,
			Grade:         r.Grade,
			Stars:         r.Stars,
			Id:            i,
			Ascents:       r.Ascents,
			Benchmark:     r.Benchmark,
			MoonId:        r.MoonId,
			RawDate:       date,
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

	prefix := strings.Replace(strings.ToLower(set.Name), " ", "", -1)
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
