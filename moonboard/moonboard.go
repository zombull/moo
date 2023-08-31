package moonboard

import (
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
	"github.com/zombull/moo/moonboard/mb1"
	"github.com/zombull/moo/moonboard/mb2"
)

const Name = "Moonboard"

func Init(d *database.Database) {
	for _, yyyy := range []string{ "2016", "2017", "2019" } {
		name := "MoonBoard " + yyyy

		if !d.Exists(name, "sets") {
			s := database.Set{
				Name:   name,
			}
			d.Insert(&s)
		}
	}
}

func SetId(d *database.Database, setName string) int64 {
	return d.GetSetId(setName)
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

func syncProblemsJSON(d *database.Database, setYear string, problems []mb.Problem) {
	setId := SetId(d, "MoonBoard " + setYear)

	for _, p := range problems {
		if p.ApiId == 0 {
			p.ApiId = p.Id
			bug.On(p.ApiId == 0, fmt.Sprintf("Moonboard problem '%s' ApiId and ID are both zero :(", p.Name))
		}
		name := sanitize(p.Name)
		orig := d.FindProblem(setId, name)
		if orig != nil {
			// Moonboard apparently has an insertion bug of some form and
			// allows for back-to-back insertions with the same name.
			if orig.MoonId != p.ApiId {
				fmt.Printf("WARN: duplicate problem found: %s %d -> %d\n", name, orig.MoonId, p.ApiId)
				continue
			}
		} else {
			orig = d.FindProblemByMoonId(setId, p.ApiId)
			if orig != nil {
				fmt.Printf("BUG: Moonboard problem exists: \"%s\": { Name: \"%s\", ApiId: %d},\n", p.Name, orig.Name, p.ApiId)
				continue
			}
		}

		setter := d.FindSetter(p.Setter.Name)
		if setter == nil {
			if orig != nil {
				setter = d.GetSetter(orig.SetterId)
				// fmt.Printf("WARN: Moonboard problem '%s' using setter '%s' instead of '%s'\n", name, setter.Name, p.Setter.Name)
			} else {
				setter = &database.Setter{
					Name:     p.Setter.Name,
				}
				if p.Setter.Name != sanitize(p.Setter.Nickname) {
					setter.Nickname = p.Setter.Nickname
				}
				d.Insert(setter)
			}
		} else if setter.Nickname != p.Setter.Nickname && len(setter.Nickname) == 0 {
			setter.Nickname = p.Setter.Nickname
			d.Update(setter)
		}

		problem := &database.Problem{
			SetId:     setId,
			SetterId:  setter.Id,
			Name:      name,
			MoonId:    p.ApiId,
			Ascents:   p.Ascents,
			Benchmark: p.Benchmark,
			Stars:     p.Stars,
		}

		var ok bool
		if !problem.Benchmark && len(p.UserGrade) > 0 {
			problem.Grade, ok = database.FontainebleauToHueco[strings.ToUpper(p.UserGrade)]
		} else {
			problem.Grade, ok = database.FontainebleauToHueco[strings.ToUpper(p.Grade)]
		}
		bug.On(!ok, fmt.Sprintf("Unhandled case in 'Grade rating': %v", p.UserGrade))

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
				if t != to {
					fmt.Printf("BUG: Duplicate hold '%s' of different type in problem '%s'\n", loc, name)
					continue
				}
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

		problem.Date = p.Date.Unix()

		if orig != nil {
			h2 := d.GetHolds(orig.Id)
			bug.On(len(h2.Holds) != len(holds.Holds), fmt.Sprintf("Existing Moonboard problem '%s' holds diverge", problem.Name))
			for i := range h2.Holds {
				bug.On(h2.Holds[i] != holds.Holds[i], fmt.Sprintf("Existing Moonboard problem '%s' holds diverge", problem.Name))
			}

			if orig.SetterId < problem.SetterId {
				problem.SetterId = orig.SetterId
			}

			bug.On(orig.MoonId != problem.MoonId, fmt.Sprintf("Existing Moonboard problem '%s' ID diverges", problem.Name))

			var o_date time.Time = time.Unix(orig.Date, 0)
			bug.On(p.Date.Format("2006-01-02") != o_date.Format("2006-01-02"), fmt.Sprintf("Existing Moonboard problem '%s' date diverges %s -> %s\n", problem.Name, o_date, p.Date))
	
			problem.Id = orig.Id
			d.Update(problem)
		} else {
			d.InsertDoubleLP(problem, holds)
		}
	}
}

func SyncProblemsJSONv1(d *database.Database, setYear string, data []byte) {
	problems := mb1.Problems{}
	err := json.Unmarshal(data, &problems)
	bug.OnError(err)

	for _, p := range problems.Problems {
		p.Setter.Name = fmt.Sprintf("%s %s", sanitize(p.Setter.Firstname), sanitize(p.Setter.Lastname))
	}

	bug.Bug("v1 deprecated")
	// syncProblemsJSON(d, setYear, problems.Problems)
}

type rename struct {
	Name string
	ApiId uint
}
var renames = map[string]*rename{
	"born Slippy": { Name: "Born Slippy 2", ApiId: 360798},
	"Campus": { Name: "Campus 2", ApiId: 386568},
	"TWIX": { Name: "Twix 2", ApiId: 395364},
	"NEMESIS": { Name: "Nemesis 2", ApiId: 358037},
	"Shorty": { Name: "Shorty 2", ApiId: 420807},
	"Kraftwerk": { Name: "Kraftwerk 2", ApiId: 423610},
	"Easy peasy lemon squeezy": { Name: "Easy Peasy Lemon Squeezy 2", ApiId: 427738},
	"GAZPACHO_1": { Name: "GAZPACHO_1_2", ApiId: 429209},
	"ACG6 was too hard": { Name: "ACG6 WAS TOO HARD", ApiId: 89571},
	"acg34var": { Name: "ACG34VAR", ApiId: 103050},
	"CARBON &OCK": { Name: "Carbon Cock", ApiId: 176870},
	"acg59-variante": { Name: "ACG59-VARIANTE", ApiId: 177035},
	"P LIPS": { Name: "Pussy Lips", ApiId: 196512},
	"C MOUTH": { Name: "Cock Mouth", ApiId: 196515},
	"friends": { Name: "Benmoonship", ApiId: 238494},
	"S** ON THE BEACH": { Name: "Sex On The Beach", ApiId: 260898},
	"seni sordum yıldızlara": { Name: "Seni Sordum Yildizlara", ApiId: 308948},
	"ACG23 - Easy": { Name: "ACG23 - EASY", ApiId: 309580},
}

var still_useful = map[uint]bool{
	88047: true,
	162380: true,
	165950: true,
	260791: true,
	274976: true,
	308347: true,
	308908: true,
	309200: true,
	309578: true,
	309703: true,
	309777: true,
	309927: true,
	309985: true,
	310114: true,
	310145: true,
	310803: true,
	310807: true,
	310894: true,
	310947: true,
	310959: true,
	311110: true,
	311194: true,
	311921: true,
	311940: true,
	312340: true,
	312492: true,
	312593: true,
	312767: true,
	312953: true,
	313079: true,
	313210: true,
	313284: true,
	313306: true,
	313642: true,
	313665: true,
	313726: true,
	313748: true,
	313750: true,
	313754: true,
	313770: true,
	314490: true,
	314722: true,
	314776: true,
	314923: true,
	314960: true,
	314962: true,
	315060: true,
	315095: true,
	315716: true,
	315893: true,
	316071: true,
	316072: true,
	316148: true,
	316189: true,
	316479: true,
	316714: true,
	316745: true,
	316782: true,
	317014: true,
	317247: true,
	317462: true,
	317925: true,
	318284: true,
	318672: true,
	318787: true,
	319211: true,
	319641: true,
	319714: true,
	320123: true,
	320258: true,
	320529: true,
	320770: true,
	320898: true,
	320906: true,
	320980: true,
	321143: true,
	321328: true,
	321443: true,
	321771: true,
	322221: true,
	322609: true,
	322917: true,
	323688: true,
	325338: true,
	325963: true,
	326490: true,
	327503: true,
	328204: true,
	328572: true,
	328635: true,
	329408: true,
	329937: true,
	331169: true,
	331233: true,
	331766: true,
	331801: true,
	331802: true,
	331803: true,
	332830: true,
	334842: true,
	335478: true,
	336226: true,
	338306: true,
	338570: true,
	338856: true,
	341672: true,
	341731: true,
	343877: true,
	344558: true,
	345685: true,
	358592: true,
}

func isUseless(moonId, ascents uint, benchmark bool, setter string) bool {
	if ascents > 3 || benchmark {
		return false
	}
	if _, ok := still_useful[moonId]; ok {
		return false
	}
	if setter == "Kyle Knapp" {
		return false
	}
	return true
}

func SyncProblemsJSONv2(d *database.Database, problemsData, holdsData []byte) {
	var problemsV2 []mb2.Problem
	err := json.Unmarshal(problemsData, &problemsV2)
	bug.OnError(err)

	var holdsV2 []mb2.Hold
	err = json.Unmarshal(holdsData, &holdsV2)
	bug.OnError(err)

	var problems2016 []mb.Problem
	var problems2017 []mb.Problem
	var problems2019 []mb.Problem

	problems := make([]mb.Problem, len(problemsV2))

	for i, p := range problemsV2 {
		bug.On(p.Id - 1 != uint(i),  fmt.Sprintf("Moonboard problem '%s' ID '%d' != index '%d'", p.Name, p.Id, i))

		setter := mb.Setter{}
		setter.Name = sanitize(p.Setter)
		setter.Firstname = ""
		setter.Lastname = ""
		setter.Nickname = setter.Name
		setter.Id = p.SetbyId

		if re, ok := renames[p.Name]; ok && p.ApiId == re.ApiId{
			fmt.Printf("Renamed '%s' to '%s'\n", p.Name, re.Name)
			p.Name = re.Name
		}

		ms := (p.Date / 10000) - 62135596800000

		problems[i].Date      = time.Unix(0, int64(time.Millisecond) * ms)
		problems[i].ApiId     = p.ApiId
		problems[i].Id        = p.ApiId
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

		if isUseless(p.ApiId, p.Ascents, p.Benchmark, p.Setter.Name) {
			continue
		}
	
		setupId := problemsV2[i].SetupId
		if setupId == 4 {
			problems2016 = append(problems2016, p)
		} else if setupId == 2 {
			problems2017 = append(problems2017, p)
		} else if setupId == 3 {
			problems2019 = append(problems2019, p)
		} else {
			bug.Bug(fmt.Sprintf("SetupId '%d' isn't valid", setupId))
		}
	}

	syncProblemsJSON(d, "2016", problems2016)
	syncProblemsJSON(d, "2017", problems2017)
	syncProblemsJSON(d, "2019", problems2019)
}



func purge(d *database.Database, setYear string) {
	problems := d.GetProblems(SetId(d, "MoonBoard " + setYear))

	for _, p := range problems {
		if isUseless(p.MoonId, p.Ascents, p.Benchmark, d.GetSetter(p.SetterId).Name) {
			fmt.Printf("DELETE %s\n", p.Name)
			d.Delete(d.GetHolds(p.Id))
			d.Delete(p)
		}
	}
}

func Purge(d *database.Database) {
	purge(d, "2016")
	purge(d, "2017")
	purge(d, "2019")
}
