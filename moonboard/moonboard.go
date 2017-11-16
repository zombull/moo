package moonboard

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/moonboard/gd"
	"github.com/zombull/floating-castle/moonboard/mb"
)

var setName = ""

const Name = "Moonboard"

func Init(d *database.Database, set string) {
	exists := d.Exists(Name, "crags")

	setName = strings.TrimSpace(set)
	if len(setName) == 0 {
		if exists {
			fmt.Printf("Moonboard exists in your database but no set is defined in your config!\n\n")
		}
		return
	}

	if !exists {
		c := database.Crag{
			Name:     Name,
			Location: "Portland Rock Gym",
			Url:      "https://www.moonboard.com/Problems/Index",
			Map:      "https://www.moonboard.com/Problems/Index",
		}
		d.Insert(&c)
	}
	if !d.Exists(setName, "areas") {
		a := database.Area{
			CragId: d.GetCragId(Name),
			Name:   setName,
			Url:    "https://www.moonboard.com/Problems/Index",
		}
		d.Insert(&a)
	}
}

func SetDefined() bool {
	return len(setName) > 0
}

func Id(d *database.Database) int64 {
	return d.GetCragId(Name)
}

func SetId(d *database.Database) int64 {
	bug.On(!SetDefined(), "accessing undefined Moonboard set")
	return d.GetAreaId(Id(d), setName)
}

var attemptsRegex = regexp.MustCompile(`Attempts: ([0-9]+)`)
var sessionsRegex = regexp.MustCompile(`Sessions: ([0-9]+)`)
var starsRegex = regexp.MustCompile(`Stars: ([0-9]+)`)
var commentRegex = regexp.MustCompile(`Comment: ([[:ascii:]]+)`)

var triesToAttempts = map[string]uint{
	"Flashed":           1,
	"2nd try":           2,
	"3rd try":           3,
	"more than 3 tries": 10,
}

func regexFindGroup(r *regexp.Regexp, s string) (string, bool) {
	if ss := r.FindStringSubmatch(strings.TrimSpace(s)); len(ss) >= 2 {
		return strings.TrimSpace(ss[1]), true
	}
	return "", false
}

func regexParseUint(r *regexp.Regexp, s string, optional bool) (uint, bool) {
	if g, ok := regexFindGroup(r, s); ok {
		i, err := strconv.ParseUint(g, 10, 64)
		if optional {
			if err != nil || i == 0 {
				fmt.Println("Bad data in comment section, falling back to log")
				return 0, false
			}
		} else {
			bug.OnError(err)
		}

		return uint(i), true
	}
	return 0, false
}

func sanitize(s string) string {
	if strings.HasPrefix(s, "ACG") {
		return s
	}
	s = strings.Title(strings.ToLower(s))
	for k, v := range map[string]string{"'Ll": "'ll", "I'M": "I'm", "I'V": "I'v", "'S": "'s", "'T": "'t", "u'R": "u'r"} {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

var countries = map[string]string{
	"NSW":                          "australia",
	"AUSTRAILIA":                   "australia",
	"AUSTRLIA":                     "australia",
	"AUSRALIA":                     "australia",
	"AUS":                          "australia",
	"CANBERRA":                     "australia",
	"ÖSTERREICH":                   "austria",
	"AT":                           "austria",
	"DORNBIRN":                     "austria",
	"ASERBAIDSCHAN":                "azerbaijan",
	"BRASIL":                       "brazil",
	"BULGRIA":                      "bulgaria",
	"BC":                           "canada",
	"CAN":                          "canada",
	"ONTARIO":                      "canada",
	"CHAINN":                       "china",
	"CZ":                           "czech",
	"CZE":                          "czech",
	"CZECHIA":                      "czech",
	"CRO":                          "croatia",
	"HRVATSKA":                     "croatia",
	"CUNDINAMARCA":                 "colombia",
	"DDEDENDENMDENMADENMARDENMARK": "denmark",
	"DK":                        "denmark",
	"DANMARK":                   "denmark",
	"WARWICKSHIRE":              "england",
	"GREATER MANCHESTER":        "england",
	"AVON":                      "england",
	"BRISTOL":                   "england",
	"CHESHIRE":                  "england",
	"CUMBRIA":                   "england",
	"ENGLANG":                   "england",
	"GUERNSEY":                  "england",
	"NORFOLK":                   "england",
	"NORTHUMBERLAND":            "england",
	"MANCHESTER":                "england",
	"SHROPSHIRE":                "england",
	"USE":                       "england",
	"SUOMI":                     "finland",
	"D":                         "germany",
	"DE":                        "germany",
	"BREMEN":                    "germany",
	"GE":                        "germany",
	"GER":                       "germany",
	"GERMAN":                    "germany",
	"GERMENY":                   "germany",
	"GERMANNY":                  "germany",
	"DEUTSCHLAND":               "germany",
	"GERNANY":                   "germany",
	"SACHSEN":                   "germany",
	"HK":                        "hong kong",
	"CLUB":                      "italy",
	"MILANO":                    "italy",
	"IT":                        "italy",
	"ITA":                       "italy",
	"ITAY":                      "italy",
	"ITLY":                      "italy",
	"ITALIA":                    "italy",
	"VICENZA":                   "italy",
	"MEZZASELVA":                "italy",
	"CERNUSCO":                  "italy",
	"AGORDO":                    "italy",
	"CASTGNIT":                  "italy",
	"香川":                        "japan",
	"日本":                        "japan",
	"OSAKA":                     "japan",
	"JJAJAPJAPAJAPANJAPAN":      "japan",
	"JJAJAPJAPAJAPANJAPANJAPAN": "japan",
	"JAPANN":                    "japan",
	"JAP":                       "japan",
	"JPAN":                      "japan",
	"JPP":                       "japan",
	"JPN":                       "japan",
	"KYOTO":                     "japan",
	"NJAPAN":                    "japan",
	"OKINAWA":                   "japan",
	"LUX":                       "luxembourg",
	"MAROCCO":                   "morocoo",
	"HOLLAND":                   "netherlands",
	"NEDERLAND":                 "netherlands",
	"NL":                        "netherlands",
	"NERERLAND":                 "netherlands",
	"NZ":                        "new zealand",
	"OSLO":                      "norway",
	"NORGE":                     "norway",
	"PL":                        "poland",
	"POL":                       "poland",
	"POLSKA":                    "poland",
	"TIMIȘ":                     "romania",
	"ROMÂNIA":                   "romania",
	"RO":                        "romania",
	"РОССИЯ":                    "russia",
	"HAMILTON":                  "scotland",
	"COTLAND":                   "scotland",
	"한국":                        "south korea",
	"ㅋㅋㅋ":                       "south korea",
	"대한민국":                      "south korea",
	"GUNSAN":                    "south korea",
	"SOUTHKORE":                 "south korea",
	"SEOUL":                     "south korea",
	"JAA":                       "south korea",
	"KOR":                       "south korea",
	"SINJUNGDONG":               "south korea",
	"CATLUNYA":                  "spain",
	"CATALONIA":                 "spain",
	"CATALUNYA":                 "spain",
	"ASTURIAS":                  "spain",
	"BARBASTRO":                 "spain",
	"ESPAÑA":                    "spain",
	"LLIRIA":                    "spain",
	"SUPEIN":                    "spain",
	"SVERIGE":                   "sweden",
	"SWE":                       "sweden",
	"YSTAD":                     "sweden",
	"SCHWEIZ":                   "switzerland",
	"CH":                        "switzerland",
	"SUISSE":                    "switzerland",
	"TICINO":                    "switzerland",
	"SVIZZERA":                  "switzerland",
	"SCHWIZERLAND":              "switzerland",
	"VALAIS":                    "switzerland",
	"TAOYUAN":                   "taiwain",
	"UAE":                       "united arab emirates",
	"GB":                        "united kingdom",
	"GBR":                       "united kingdom",
	"BRITAIN":                   "united kingdom",
	"UK":                        "united kingdom",
	"BOULDER":                   "united states",
	"US":                        "united states",
	"HAWAII":                    "united states",
	"USA":                       "united states",
	"SAN DIEGO":                 "united states",
	"TEXAS":                     "united states",
	"TX":                        "united states",
	"HI":                        "united states",
	"MA":                        "united states",
	"NC":                        "united states",
	"SC":                        "united states",
	"CA":                        "united states",
	"CO":                        "united states",
	"COLORADO":                  "united states",
	"OR":                        "united states",
	"OREGON":                    "united states",
	"AMERICA":                   "united states",
	"MURICA":                    "united states",
	"WA":                        "united states",
	"NJ":                        "united states",
	"IN":                        "united states",
	"TN":                        "united states",
	"MD":                        "united states",
	"SANTA CLARA":               "united states",
	"LA":                        "united states",
	"MICHIGAN":                  "united states",
	"BALTIMORE":                 "united states",
	"66606":                     "united states",
	"WY":                        "united states",
	"UA":                        "united states",
	"UAUSAUSAUSA":               "united states",
	"UNITED STATES OF AMERICA":  "united states",
	"UNITED STATES":             "united states",
	"UNITEDSTATES":              "united states",
	"WASHINGTON":                "united states",
	"BRIDGEND":                  "wales",
	"CARDIFF":                   "wales",
	"GWENT":                     "wales",
	"DENBIGHSHIRE":              "wales",
	"WREXHAM":                   "wales",
	"FLINTSHIRE":                "wales",
	"TRAINYOURBICEPS":           "moon",
	"DADDY":                     "moon",
	"TBD":                       "moon",
	"CDGH":                      "moon",
	"THH":                       "moon",
	"NO":                        "moon",
	"GONDOLA":                   "moon",
	"TRF":                       "moon",
	"0988614220":                "moon",
}

func country(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	if c, ok := countries[s]; ok {
		s = c
	}
	return sanitize(s)
}

var forcedNicknames = map[string]string{
	"Jon Guinther":   "goonthorj",
	"Mark Tomlinson": "threenine",
	"Juha Isotupa":   "JuIs2000",
}

func SyncProblems(d *database.Database, data []byte) {
	cragId := Id(d)
	setId := SetId(d)

	problems := mb.Problems{}
	err := json.Unmarshal(data, &problems)
	bug.OnError(err)

	for _, p := range problems.Problems {
		name := sanitize(p.Name)
		route := d.FindRoute(setId, name)
		if route != nil {
			// Moonboard apparently has an insertion bug of some form and
			// allows for back-to-back insertions with the same name.
			if route.Length != p.ApiId {
				fmt.Printf("WARN: duplicate route found: %s\n", name)
			} else {
				// Updates not yet implemented
			}
			continue
		}

		sname := fmt.Sprintf("%s %s", sanitize(p.Setter.Firstname), sanitize(p.Setter.Lastname))
		setter := d.FindSetter(cragId, sname)
		if setter == nil {
			setter = &database.Setter{
				CragId:   cragId,
				Name:     sname,
				Country:  country(p.Setter.Country),
				City:     sanitize(p.Setter.City),
				Inactive: false,
			}
			if fname, ok := forcedNicknames[sname]; ok {
				setter.Nickname = fname
			} else if sname != sanitize(p.Setter.Nickname) {
				setter.Nickname = p.Setter.Nickname
			}
			d.Insert(setter)
		}
		if _, ok := forcedNicknames[sname]; !ok && sname != sanitize(p.Setter.Nickname) {
			bug.On(setter.Nickname != p.Setter.Nickname, fmt.Sprintf("Duplicate setter? %v vs. %v", setter, p.Setter))
		}

		route = &database.Route{
			CragId:    cragId,
			AreaId:    setId,
			Name:      name,
			Type:      "moonboard",
			Length:    p.ApiId,
			Pitches:   p.Ascents,
			Benchmark: p.Benchmark,
			Stars:     p.Stars,
		}
		route.SetterId = setter.Id

		if len(p.Url) == 0 {
			b, err := json.Marshal(p)
			bug.OnError(err)
			p.Url = fmt.Sprintf("%x", md5.Sum(b))
		}
		route.Url = fmt.Sprintf("https://www.moonboard.com/Problems/View/%d/%s", p.ApiId, p.Url)

		var ok bool
		if len(p.UserGrade) > 0 {
			route.Grade, ok = database.FontainebleauToHueco[strings.ToUpper(p.UserGrade)]
		} else {
			route.Grade, ok = database.FontainebleauToHueco[strings.ToUpper(p.Grade)]
		}

		bug.On(!ok, fmt.Sprintf("Unhandled case in 'Grade rating': %v", p.UserGrade))

		route.Date, err = time.Parse("02 Jan 2006 15:04", p.Date)
		bug.OnError(err)

		holds := &database.Holds{
			Holds: make([]string, 0, len(p.Holds)),
		}

		holdMap := make(map[string]string)
		for _, h := range p.Holds {
			var loc = strings.ToUpper(h.Location)
			t := "i"
			if h.IsStart {
				t = "s"
			} else if h.IsEnd {
				t = "f"
			}
			if to, ok := holdMap[loc]; ok {
				bug.On(t != to, fmt.Sprintf("Duplicate hold '%s' of different type in problem '%s'", loc, name))
			} else {
				holds.Holds = append(holds.Holds, t+loc)
				holdMap[loc] = t
			}
		}
		sort.Strings(holds.Holds)

		d.InsertDoubleLP(route, holds)
	}
}

var attempts = map[string]uint{
	"ACG23":                   25,
	"Mark's Favorite Problem": 50,
}

var sessions = map[string]uint{
	"ACG23":                   5,
	"Mark's Favorite Problem": 10,
}

var stars = map[string]uint{
	"Saltedblocks 2":          5,
	"Tca Comp 1":              5,
	"Left Jab":                5,
	"Hard Times":              5,
	"Kang Mina":               5,
	"ACG24":                   5,
	"46":                      5,
	"Intimissimi #2":          5,
	"ACG58":                   5,
	"The Famous Five":         5,
	"Problem 9":               5,
	"Air Force":               5,
	"ACG23":                   5,
	"Kicker":                  5,
	"Wills Var":               5,
	"Austroraptor":            5,
	"2 Hours Of Purity Ring":  5,
	"5 Finger Discount":       5,
	"Mark's Favorite Problem": 5,
}

func SyncTicks(d *database.Database, data []byte) {
	cragId := Id(d)
	setId := SetId(d)

	ticks := mb.Ticks{}
	err := json.Unmarshal(data, &ticks)
	bug.OnError(err)

	for _, t := range ticks.Ticks {
		route := d.FindRoute(setId, sanitize(t.Problem.Name))
		bug.On(route == nil, fmt.Sprintf("'%s' not in database, Moonboard index needs to be synced\n", t.Problem.Name))

		if e := d.GetTicks(route.Id); len(e) > 0 {
			continue
		}

		tick := &database.Tick{
			RouteId:  route.Id,
			AreaId:   setId,
			CragId:   cragId,
			Stars:    t.Stars,
			Redpoint: true,
		}

		if ui, ok := attempts[route.Name]; ok {
			tick.Attempts = ui
		} else if ui, ok := regexParseUint(attemptsRegex, t.Comment, true); ok {
			tick.Attempts = ui
		} else {
			tick.Attempts, ok = triesToAttempts[t.NumberOfTries]
			bug.On(!ok, fmt.Sprintf("Unhandled case in 'Number of tries': %s", t.NumberOfTries))

		}
		if ui, ok := sessions[route.Name]; ok {
			tick.Sessions = ui
		} else if ui, ok := regexParseUint(sessionsRegex, t.Comment, true); ok {
			tick.Sessions = ui
		}
		if ui, ok := stars[route.Name]; ok {
			tick.Stars = ui
		} else if ui, ok := regexParseUint(starsRegex, t.Comment, true); ok {
			tick.Stars = ui
		}
		if s, ok := regexFindGroup(commentRegex, t.Comment); ok {
			tick.Comment = s
		}

		var ok bool
		tick.Grade, ok = database.FontainebleauToHueco[strings.ToUpper(t.Grade)]
		bug.On(!ok, fmt.Sprintf("Unhandled case in 'Grade rating': %s", t.Grade))

		tick.Date, err = time.Parse("02 Jan 2006", t.Date)
		bug.OnError(err)

		tick.Flash = (tick.Attempts == 1)
		d.Insert(tick)
	}
}

func SyncUserData(d *database.Database, u *gd.UserData) {
	ticks := make([]*database.Tick, 0, len(u.Ticks))
	cragId := Id(d)
	setId := SetId(d)
	for problem, t := range u.Ticks {
		route := d.FindRoute(setId, problem)
		bug.On(route == nil, fmt.Sprintf("'%s' not in database, Moonboard index needs to be synced", problem))

		if existing := d.GetTicks(route.Id); len(existing) > 0 {
			continue
		}

		_, ok := database.HuecoToFontainebleau[strings.ToUpper(t.Grade)]
		bug.On(!ok, fmt.Sprintf("Invalid grade '%s' for problem '%s'", t.Grade, problem))

		tick := &database.Tick{
			RouteId:  route.Id,
			AreaId:   setId,
			CragId:   cragId,
			Grade:    t.Grade,
			Stars:    t.Stars,
			Attempts: t.Attempts,
			Flash:    t.Attempts == 1,
			Redpoint: true,
		}
		var err error
		tick.Date, err = time.Parse("January 02, 2006", t.Date)
		bug.OnError(err)

		if t.Sessions > 0 {
			tick.Sessions = t.Sessions
		}
		ticks = append(ticks, tick)
	}
	for _, tick := range ticks {
		d.Insert(tick)
	}
}
