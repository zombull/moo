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

	"github.com/zombull/moo/bug"
	"github.com/zombull/moo/database"
	"github.com/zombull/moo/moonboard/mb"
	"github.com/zombull/moo/moonboard/mb2"
)

const Name = "Moonboard"

func Init(d *database.Database) {
	exists := d.Exists(Name, "crags")
	if !exists {
		c := database.Crag{
			Name:     Name,
			Location: "Portland Rock Gym",
			Url:      "https://www.moonboard.com",
			Map:      "https://www.moonboard.com",
		}
		d.Insert(&c)
	}

	cragId := d.GetCragId(Name)

	for _, yyyy := range []string{ "2016", "2017", "2019" } {
		s := "MoonBoard " + yyyy

		if !d.Exists(s, "areas") {
			a := database.Area{
				CragId: cragId,
				Name:   s,
				Url:    "https://www.moonboard.com/Problems/Index",
			}
			d.Insert(&a)
		}
	}
}

func CragId(d *database.Database) int64 {
	return d.GetCragId(Name)
}

func SetId(d *database.Database, setName string) int64 {
	return d.GetAreaId(CragId(d), setName)
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
	s = strings.TrimSpace(s)

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

func syncProblemsJSON(d *database.Database, setYear string, problems []mb.Problem) {
	cragId := CragId(d)
	setId := SetId(d, "MoonBoard " + setYear)

	var err error

	for _, p := range problems {
		if p.ApiId == 0 {
			p.ApiId = p.Id
			bug.On(p.ApiId == 0, fmt.Sprintf("Moonboard problem '%s' ApiId and ID are both zero :(", p.Name))
		}
		name := sanitize(p.Name)
		orig := d.FindRoute(setId, name)
		if orig != nil {
			// Moonboard apparently has an insertion bug of some form and
			// allows for back-to-back insertions with the same name.
			if orig.Length != p.ApiId {
				fmt.Printf("WARN: duplicate route found: %s\n", name)
				continue
			}
		} else {
			orig = d.FindRouteByLength(setId, p.ApiId)
			if orig != nil {
				bug.Bug(fmt.Sprintf("Moonboard problem '%s' with  ID '%d' exists as '%s'", name, p.ApiId, orig.Name))
			}
		}

		setter := d.FindSetter(cragId, p.Setter.Name)
		if setter == nil {
			setter = &database.Setter{
				CragId:   cragId,
				Name:     p.Setter.Name,
				Country:  country(p.Setter.Country),
				City:     sanitize(p.Setter.City),
				Inactive: false,
			}
			if p.Setter.Name != sanitize(p.Setter.Nickname) {
				setter.Nickname = p.Setter.Nickname
			}
			d.Insert(setter)
		} else if setter.Nickname != p.Setter.Nickname && len(setter.Nickname) == 0 {
			setter.Nickname = p.Setter.Nickname
			d.Update(setter)
		}

		route := &database.Route{
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
			route.Url = fmt.Sprintf("https://www.moonboard.com/Problems/View/%d/%x", p.ApiId, md5.Sum(b))
		} else {
			route.Url = fmt.Sprintf("https://www.moonboard.com/Problems/View/%d/%s", p.ApiId, p.Url)
			if orig != nil {
				bug.On(orig.Url != route.Url, fmt.Sprintf("Existing Moonboard problem '%s' url diverges: '%s' -> '%s'", route.Name, orig.Url, route.Url))
			}
		}

		var ok bool
		if !route.Benchmark && len(p.UserGrade) > 0 {
			route.Grade, ok = database.FontainebleauToHueco[strings.ToUpper(p.UserGrade)]
		} else {
			route.Grade, ok = database.FontainebleauToHueco[strings.ToUpper(p.Grade)]
			if orig != nil && orig.Grade != route.Grade{
				fmt.Printf("WARN: problem '%s' grade changed from '%s' to '%s'\n", route.Name, orig.Grade, route.Grade)
			}
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
		if len(holds.Holds) > 24 {
			fmt.Printf("WARN: skipping '%s' as it has %d holds\n", name, len(holds.Holds))
			continue
		}
		sort.Strings(holds.Holds)

		if orig != nil {
			h2 := d.GetHolds(orig.Id)
			bug.On(len(h2.Holds) != len(holds.Holds), fmt.Sprintf("Existing Moonboard problem '%s' holds diverge", route.Name))
			for i := range h2.Holds {
				bug.On(h2.Holds[i] != holds.Holds[i], fmt.Sprintf("Existing Moonboard problem '%s' holds diverge", route.Name))
			}
			bug.On(orig.SetterId != route.SetterId, fmt.Sprintf("Existing Moonboard problem '%s' setter ID diverges", route.Name))
			bug.On(orig.Date != route.Date, fmt.Sprintf("Existing Moonboard problem '%s' date diverges", route.Name))
			bug.On(orig.Length != route.Length, fmt.Sprintf("Existing Moonboard problem '%s' ID diverges", route.Name))

			route.Id = orig.Id
			d.Update(route)
		} else {
			d.InsertDoubleLP(route, holds)
		}
	}
}

func SyncProblemsJSONv1(d *database.Database, setYear string, data []byte) {
	problems := mb.Problems{}
	err := json.Unmarshal(data, &problems)
	bug.OnError(err)

	for _, p := range problems.Problems {
		p.Setter.Name = fmt.Sprintf("%s %s", sanitize(p.Setter.Firstname), sanitize(p.Setter.Lastname))
	}

	syncProblemsJSON(d, setYear, problems.Problems)
}

type rename struct {
	Name string
	Setter string
}
var renames = map[string]*rename{
	"born Slippy": { Name: "Born Slippy 2", Setter: "4e78bd45-9445-4cfb-ade7-4f91ac6d3714"},
	"Campus": { Name: "Campus 2", Setter: "66de8d9d-4835-4933-ac32-70813df4ca8a"},
	"TWIX": { Name: "Twix 2", Setter: "9f3b703a-13f5-4c35-95dd-2611b65c5e9f"},
	"NEMESIS": { Name: "Nemesis 2", Setter: "7C9A8746-7ED8-4657-9DE5-CAE2D05405C7"},
	"Shorty": { Name: "Shorty 2", Setter: "3673740f-6bd4-4995-94ad-57e373554f8c"},
	"Kraftwerk": { Name: "Kraftwerk 2", Setter: "b908ed44-3a6a-4dec-aa12-7ef8a35a4229"},
	"Easy peasy lemon squeezy": { Name: "Easy Peasy Lemon Squeezy 2", Setter: "77a57172-ecb2-4798-b0a6-0c9b43acb67a"},
	"GAZPACHO_1": { Name: "GAZPACHO_1_2", Setter: "e59e5a76-0d91-4785-befa-51ece8a15a49"},
}

func SyncProblemsJSONv2(d *database.Database, setYear string, problemsData, holdsData []byte) {
	var problemsV2 []mb2.Problem
	err := json.Unmarshal(problemsData, &problemsV2)
	bug.OnError(err)

	var holdsV2 []mb2.Hold
	err = json.Unmarshal(holdsData, &holdsV2)
	bug.OnError(err)

	problems := make([]mb.Problem, len(problemsV2))

	for i, p := range problemsV2 {
		bug.On(p.Id - 1 != uint(i),  fmt.Sprintf("Moonboard problem '%s' ID '%d' != index '%d'", p.Name, p.Id, i))

		ms := (p.Date / 10000) - 62135596800000
		date := time.Unix(0, int64(time.Millisecond) * ms)

		setter := mb.Setter{}
		setter.Country = "Unknown"
		setter.City =    "Unknown"
		setter.Name = sanitize(p.Setter)
		setter.Firstname = ""
		setter.Lastname = ""
		setter.Nickname = setter.Name
		setter.Id = p.SetbyId

		if re, ok := renames[p.Name]; ok && p.SetbyId == re.Setter {
			fmt.Printf("Renamed '%s' to '%s'\n", p.Name, re.Name)
			p.Name = re.Name
		}

		problems[i].Date      = date.Format("02 Jan 2006 15:04")
		problems[i].ApiId     = p.ApiId
		problems[i].Id        = p.ApiId
		problems[i].Url       = strings.ToLower(p.Name)
		problems[i].Setter    = setter
		problems[i].UserGrade = p.UserGrade
		problems[i].Grade     = p.Grade
		problems[i].Name      = p.Name
		problems[i].Rating    = 0
		problems[i].Stars     = p.Stars
		problems[i].Ascents   = p.Ascents
		problems[i].Benchmark = p.IsBenchmark != 0
		problems[i].Holds     = make([]mb.Hold, 0)
	}

	for _, h2 := range holdsV2 {
		idx := h2.ProblemId - 1
		bug.On(idx > uint(len(problems)), "Hold problem ID out of range")

		h := mb.Hold{}
		h.Id = 0
		h.Location = h2.Location
		h.IsStart = h2.IsStart != 0
		h.IsEnd = h2.IsEnd != 0
		problems[idx].Holds = append(problems[idx].Holds, h)
	}

	for i, p := range problems {
		bug.On(p.ApiId == 0 || len(p.Name) == 0, fmt.Sprintf("Problem '%d' is empty", i))
	}

	syncProblemsJSON(d, setYear, problems)
}

func Transfer(d, src *database.Database) {
	crags := src.GetCrags()
	bug.On(len(crags) != 1, "Multiple crags in the source, I'm lazy...")

	areas := src.GetAreas(crags[0].Id)
	bug.On(len(areas) != 1, "Multiple areas in the source, I'm lazy...")

	routes := src.GetRoutes(crags[0].Id, areas[0].Id)
	bug.On(len(areas) == 0, "No routes found in the source")

	cragId := CragId(d)
	setId := SetId(d, areas[0].Name)

	for _, r := range routes {

		h := src.GetHolds(r.Id)
		s := src.GetSetter(r.SetterId)

		h.RouteId = 0
		r.Id = 0
		r.CragId = cragId
		r.AreaId = setId

		setter := d.FindSetter(cragId, s.Name)
		if setter == nil {
			setter = s
			d.Insert(setter)
		} else if setter.Nickname != s.Nickname && len(setter.Nickname) == 0 {
			setter.Nickname = s.Nickname
			d.Update(setter)
		}
		r.SetterId = setter.Id

		d.InsertDoubleLP(r, h)
	}
}
