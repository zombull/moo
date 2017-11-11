package interactive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/moonboard"
)

type op interface {
	name() string
	crag(d *database.Database, line string)
	area(d *database.Database, line string)
	setter(d *database.Database)
	route(d *database.Database)
	tick(d *database.Database)
	list(d *database.Database)
}

var ops = map[string]op{
	"add":    &addOp{},
	"edit":   &editOp{},
	"search": &searchOp{},
	"query":  &queryOp{},
	"stats":  &statsOp{},
}

func Run(d *database.Database) {
	m := ops
	m["import"] = nil
	m["export"] = nil
	ac := newMapAutocompleter(m)
	l := newReader("Select Action: ", ac)
	doReadline(l, true, func(line string) string {
		if o, ok := ops[line]; ok && o != nil {
			runL1(o, d)
			fmt.Println()
		} else if line == "import" {
			import_(d)
		} else if line == "export" {
			export(d)
		} else {
			fmt.Println("Invalid Action: " + line)
		}
		return ""
	})
}

func getSet(s database.Set, name string) string {
	ac := newSetAutocompleter(s)
	l := newReader(name+": ", ac)
	return doReadline(l, false, func(line string) string {
		if _, ok := s[line]; !ok {
			fmt.Printf("Invalid %s: %s\n", name, line)
			return ""
		}
		return line
	})
}

func runL1(o op, d *database.Database) {
	t := getSet(database.TableTypes, o.name()+" Type")
	switch t {
	case "crag", "gym":
		o.crag(d, t)
	case "area", "wall":
		o.area(d, t)
	case "setter":
		o.setter(d)
	case "route":
		o.route(d)
	case "tick":
		o.tick(d)
	case "list":
		o.list(d)
	default:
		bug.Bug(fmt.Sprintf("unhandled type: %s", t))
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// Block search features, we don't use a history because
	// things are too contextual and get too fragmented.
	case readline.CharBckSearch, readline.CharFwdSearch:
		return r, false
	}
	return r, true
}

func newReader(prompt string, ac readline.AutoCompleter) *readline.Instance {
	c := readline.Config{
		Prompt:              prompt,
		AutoComplete:        ac,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		FuncFilterInputRune: filterInput,
	}

	l, err := readline.NewEx(&c)
	bug.OnError(err)
	return l
}

func doReadline(l *readline.Instance, catch bool, process func(string) string) string {
	if catch {
		defer func() {
			if r := recover(); r != nil {
				panic(r)
			}
		}()
	}

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				panic(nil)
			}
			continue
		} else if err == io.EOF {
			bug.OnError(err)
		}
		line = strings.TrimSpace(line)
		if process(line) != "" {
			return line
		}
	}
}

func getBool(p string) bool {
	ac := readline.NewPrefixCompleter(
		readline.PcItem("yes"),
		readline.PcItem("no"),
	)

	l := newReader(p+": ", ac)
	s := doReadline(l, false, func(line string) string {
		if line != "yes" && line != "no" {
			fmt.Printf("%s must be 'yes' or 'no'\n", p)
			return ""
		}
		return line
	})
	return s == "yes"
}

func getUint(p string) uint {
	val := uint(1)
	l := newReader(p+": ", nil)
	doReadline(l, false, func(line string) string {
		if i, err := strconv.Atoi(line); err == nil && i > 0 {
			val = uint(i)
			return line
		}
		fmt.Printf("%s must be > 0\n", p)
		return ""
	})
	return val
}

func getString(name string, allowEmpty bool) string {
	s := ""
	l := newReader(name+": ", nil)
	doReadline(l, false, func(line string) string {
		if len(line) == 0 && !allowEmpty {
			fmt.Printf("%s cannot be empty\n", name)
			return ""
		}
		s = line
		return "done"
	})
	return s
}

func getComment() string {
	return getString("Comment", true)
}

func getVGrade() string {
	return getSet(database.Hueco, "V Grade")
}

func getYdsGrade() string {
	return getSet(database.YDS, "YDS Grade")
}

func getType(d *database.Database) string {
	return getSet(database.RouteTypes, "Route Type")
}

type FilePrefixCompleter struct {
}

func (f *FilePrefixCompleter) Print(prefix string, level int, buf *bytes.Buffer) {

}

func (f *FilePrefixCompleter) GetName() []rune {
	return []rune{}
}

func (f *FilePrefixCompleter) IsDynamic() bool {
	return true
}

func (f *FilePrefixCompleter) GetChildren() []readline.PrefixCompleterInterface {
	return nil
}

func (f *FilePrefixCompleter) SetChildren(children []readline.PrefixCompleterInterface) {

}
func (f *FilePrefixCompleter) Do(line []rune, pos int) (newLine [][]rune, offset int) {
	return [][]rune{}, 0
}

func (f *FilePrefixCompleter) GetDynamicNames(line []rune) [][]rune {
	dir := path.Dir(strings.TrimSpace(string(line)))
	if infos, err := ioutil.ReadDir(dir); err == nil {
		names := make([][]rune, len(infos))
		for i, fi := range infos {
			if fi.IsDir() {
				names[i] = []rune(path.Join(dir, fi.Name()) + "/")
			} else if fi.Mode().IsRegular() {
				names[i] = []rune(path.Join(dir, fi.Name()))
			}
		}
		return names
	}
	return [][]rune{}
}

func getFiles(name string) []string {
	ac := readline.NewPrefixCompleter(&FilePrefixCompleter{})
	l := newReader(strings.Title(name)+" File(s): ", ac)

	var files []string
	doReadline(l, false, func(line string) string {
		if strings.Contains(line, "*") {
			dir := path.Dir(strings.TrimSpace(string(line)))
			infos, err := ioutil.ReadDir(dir)
			if err != nil {
				fmt.Printf("Wildcard used but directory is invalid: %s\n", line)
				return ""
			}
			for _, fi := range infos {
				if fi.Mode().IsRegular() {
					name := path.Join(dir, fi.Name())
					ok, err := path.Match(line, name)
					if err != nil {
						fmt.Printf("Invalid glob pattern: %s\n", line)
						return ""
					} else if ok {
						files = append(files, name)
					}
				}
			}
			if len(files) == 0 {
				fmt.Printf("No files match the glob pattern: %s\n", line)
				return ""
			}
			return line
		} else if fi, err := os.Stat(line); err != nil {
			fmt.Printf("%s: %s\n", line, err.Error())
			return ""
		} else if !fi.Mode().IsRegular() {
			fmt.Printf("'%s' is not a file\n", line)
			return ""
		}
		files = []string{line}
		return line
	})
	return files
}

func getInt(prompt string, m map[string]int) int {
	val := 0
	ac := newMapAutocompleter(m)
	l := newReader(prompt+": ", ac)
	doReadline(l, false, func(line string) string {
		var ok bool
		if val, ok = m[line]; !ok {
			fmt.Printf("Invalid %s: %s\n", prompt, line)
			return ""
		}
		return line
	})
	return val
}

func getStars() uint {
	return uint(getInt("Stars", database.Stars))
}

func getYear() int {
	y := time.Now().Year()
	bug.UserBugOn(y < 2017, "your system clock is wrong, update it")
	m := makeMapInts(2000, y-1999)
	return getInt("Year", m)
}

var months = map[string]time.Month{
	"January":   time.January,
	"February":  time.February,
	"March":     time.March,
	"April":     time.April,
	"May":       time.May,
	"June":      time.June,
	"July":      time.July,
	"August":    time.August,
	"September": time.September,
	"October":   time.October,
	"November":  time.November,
	"December":  time.December,
}

var daysInMonths = map[time.Month]int{
	time.January:   31,
	time.February:  28,
	time.March:     31,
	time.April:     30,
	time.May:       31,
	time.June:      30,
	time.July:      31,
	time.August:    31,
	time.September: 30,
	time.October:   31,
	time.November:  30,
	time.December:  31,
}

func getDaysInMonth(month time.Month, year int) int {
	days := daysInMonths[month]
	if month == time.February && ((year-2000)%4) == 0 {
		days++
	}
	return days
}

func getMonth() time.Month {
	ac := newMapAutocompleter(months)
	l := newReader("Month: ", ac)
	m := doReadline(l, false, func(line string) string {
		if _, ok := months[line]; !ok {
			fmt.Println("Invalid Month: " + line)
			return ""
		}
		return line
	})
	return months[m]
}

func getDay(month time.Month, year int) int {
	m := makeMapInts(1, getDaysInMonth(month, year))
	return getInt("Day", m)
}

func getDate() time.Time {
	var date time.Time

	ac := readline.NewPrefixCompleter(
		readline.PcItem("today"),
		readline.PcItem("yesterday"),
		readline.PcItem("other"),
	)

	l := newReader("Date: ", ac)
	doReadline(l, false, func(line string) string {
		switch line {
		case "today":
			date = time.Now()
		case "yesterday":
			date = time.Now().AddDate(0, 0, -1)
		case "other":
			year := getYear()
			month := getMonth()
			day := getDay(month, year)
			date = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		default:
			fmt.Println("Invalid Date: " + line)
			return ""
		}
		return line
	})
	return date
}

func getUrl(name string) string {
retry:
	s := getString(name, true)
	if len(s) > 0 {
		if _, err := url.ParseRequestURI(s); err != nil {
			fmt.Printf("'%s' is not a valid URL\n", s)
			goto retry
		}
	}
	return s

}

func getName(m map[string]interface{}, name string) string {
	ac := newMapAutocompleter(m)
	l := newReader(name+" Name: ", ac)
	return doReadline(l, false, func(line string) string {
		var ok bool
		if _, ok = m[line]; ok {
			return line
		}
		fmt.Printf("Select a valid %s name!\n", name)
		return ""
	})
}

func getCrag(d *database.Database) *database.Crag {
	crags := d.GetCrags()
	bug.UserBugOn(len(crags) == 0, "add a crag first")
	if len(crags) == 1 {
		return crags[0]
	}

	m := makeMapCrags(crags)
	name := getName(m, "Crag/Gym")
	return m[name].(*database.Crag)
}

func getArea(d *database.Database, crag *database.Crag) *database.Area {
	if crag.Name == moonboard.Name && moonboard.SetDefined() {
		return d.GetArea(moonboard.SetId(d))
	}
	areas := d.GetAreas(crag.Id)
	bug.UserBugOn(len(areas) == 0, fmt.Sprintf("add an area to %s first", crag.Name))
	if len(areas) == 1 {
		return areas[0]
	}

	m := makeMapAreas(areas)
	name := getName(m, "Area/Wall")
	return m[name].(*database.Area)
}

func getRoute(d *database.Database, crag *database.Crag, area *database.Area) *database.Route {
	if area == nil && crag.Name == moonboard.Name && moonboard.SetDefined() {
		area = d.GetArea(moonboard.SetId(d))
	}
	var routes []*database.Route
	if area != nil {
		routes = d.GetRoutes(crag.Id, area.Id)
	} else {
		routes = d.GetAllRoutes(crag.Id)
	}
	if len(routes) == 0 {
		return nil
	}

	m := makeMapRoutes(routes)
	m["*new*"] = nil
	if area == nil {
		m["*area*"] = nil
	}
	name := getName(m, "Route")
	if name == "*new*" {
		return nil
	} else if name == "*area*" {
		return getRoute(d, crag, getArea(d, crag))
	}
	return m[name].(*database.Route)
}

func getSetter(d *database.Database, crag *database.Crag) *database.Setter {
	setters := d.GetSetters(crag.Id)
	if len(setters) == 0 {
		return nil
	}

	m := makeMapSetters(setters)
	name := getName(m, "Setter")
	return m[name].(*database.Setter)
}
