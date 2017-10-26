package customs

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/database"
	"gopkg.in/yaml.v2"
)

type Gym struct {
	Gym      string   `yaml:"gym"`
	Location string   `yaml:"location,omitempty"`
	Url      string   `yaml:"url,omitempty"`
	Map      string   `yaml:"map,omitempty"`
	Comment  string   `yaml:"comment,omitempty"`
	Walls    []Wall   `yaml:"walls"`
	Setters  []Setter `yaml:"setters"`
}

type Wall struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url,omitempty"`
	Map     string `yaml:"map,omitempty"`
	Comment string `yaml:"comment,omitempty"`
}

type Setter struct {
	Name     string `yaml:"name"`
	Inactive bool   `yaml:"inactive,omitempty"`
	Comment  string `yaml:"comment,omitempty"`
}

func validateURL(s string) string {
	_, err := url.ParseRequestURI(s)
	bug.UserBugOn(err != nil, fmt.Sprintf("'%s' is not a valid URL\n", s))
	return s
}

func ImportGym(d *database.Database, files []string) {
	for _, file := range files {
		fmt.Printf("Importing gym from '%s'\n", file)
		importGym(d, file)
	}
}

func importGym(d *database.Database, path string) {
	data, err := ioutil.ReadFile(path)
	bug.UserError(err)

	var g Gym
	err = yaml.Unmarshal(data, &g)
	bug.UserError(err)

	bug.UserBugOn(d.Exists(g.Gym, "crags"), fmt.Sprintf("the gym '%s' already exists in your database", g.Gym))

	gym := &database.Crag{
		Name:     g.Gym,
		Location: g.Location,
		Url:      validateURL(g.Url),
		Map:      validateURL(g.Map),
		Comment:  g.Comment,
	}

	x := len(g.Walls)
	records := make([]database.SideTwo, x+len(g.Setters))
	for i, w := range g.Walls {
		records[i] = &database.Area{
			Name:    w.Name,
			Url:     w.Url,
			Map:     w.Map,
			Comment: w.Comment,
		}
	}

	for i, s := range g.Setters {
		records[x+i] = &database.Setter{
			Name:     s.Name,
			Inactive: s.Inactive,
			Comment:  s.Comment,
		}
	}

	d.InsertCollection(gym, records)
}

type GymSet struct {
	Gym   string    `yaml:"gym"`
	Date  string    `yaml:"date"`
	Walls []GymWall `yaml:"walls"`
}

type GymWall struct {
	Name   string     `yaml:"name"`
	Height uint       `yaml:"height,omitempty"`
	Routes []GymRoute `yaml:"routes"`
}

type GymRoute struct {
	Color    string `yaml:"color"`
	Grade    string `yaml:"grade"`
	Setter   string `yaml:"setter"`
	TopRope  bool   `yaml:"toprope,omitempty"`
	Redpoint bool   `yaml:"redpoint"`
	Flash    bool   `yaml:"flash,omitempty"`
	Onsight  bool   `yaml:"onsight,omitempty"`
	Attempts uint   `yaml:"attempts"`
	Sessions uint   `yaml:"sessions"`
	Falls    uint   `yaml:"falls,omitempty"`
	Hangs    uint   `yaml:"hangs,omitempty"`
	Stars    uint   `yaml:"stars"`
	Comment  string `yaml:"comment,omitempty"`
}

func default1(v uint) uint {
	if v == 0 {
		return 1
	}
	return v
}

func getStars(route string, stars uint) uint {
	bug.UserBugOn(stars < 1 || stars > 5, fmt.Sprintf("%s: stars must be between 1 and 5", route))
	return stars
}

func getGymGrade(grade string) (string, string) {
	if _, ok := database.Hueco[grade]; ok {
		return grade, "boulder"
	} else if _, ok := database.YDS[grade]; ok {
		return grade, "sport"
	}
	bug.UserBug(fmt.Sprintf("'%s' is not a valid gym grade, must be a V or YDS grade", grade))
	return "", ""
}

func getGymLead(tr bool, rtype string) bool {
	bug.UserBugOn(tr && rtype != "sport", "toprope cannot be true for Gym boulder problems")
	return !tr
}

func getRFO(attempts, falls, hangs uint, flash, onsight bool) (bool, bool, bool) {
	bug.UserBugOn((falls > 0 || hangs > 0 || attempts > 1) && (flash || onsight), "flash/onsight cannot be true when falls/hangs > 0 or attempts > 1")
	bug.UserBugOn(flash && onsight, "you can't flash and onsight a problem, pick one")
	redpoint := flash || onsight || (falls == 0 && hangs == 0)
	flash = flash || (!onsight && redpoint && attempts == 1)
	return redpoint, flash, onsight
}

type RouteAndTick struct {
	route *database.Route
	tick  *database.Tick
}

func ImportGymSet(d *database.Database, files []string) {
	for _, file := range files {
		fmt.Printf("Importing set from '%s'\n", file)
		importGymSet(d, file)
	}
}

func importGymSet(d *database.Database, path string) {
	routes := make(map[string]RouteAndTick)

	data, err := ioutil.ReadFile(path)
	bug.UserError(err)

	var gs GymSet
	err = yaml.Unmarshal(data, &gs)
	bug.UserError(err)

	gym := d.FindCrag(gs.Gym)
	bug.UserBugOn(gym == nil, fmt.Sprintf("the gym '%s' does not exist in your database", gs.Gym))

	date, err := time.Parse("2006-01-02", gs.Date)
	bug.UserError(err)

	for _, w := range gs.Walls {
		wall := d.FindArea(gym.Id, w.Name)
		bug.UserBugOn(wall == nil, fmt.Sprintf("the wall '%s' does not exist in your database for %s", w.Name, gs.Gym))

		for _, r := range w.Routes {
			setter := d.FindSetter(gym.Id, r.Setter)
			bug.UserBugOn(setter == nil, fmt.Sprintf("the setter '%s' does not exist in your database for %s", r.Setter, gs.Gym))

			route := &database.Route{
				CragId:   gym.Id,
				AreaId:   wall.Id,
				SetterId: setter.Id,
				Length:   w.Height,
				Pitches:  1,
			}

			route.Name = fmt.Sprintf("%s %s %s", gs.Date, r.Color, r.Grade)
			if _, exists := routes[route.Name]; exists {
				for i := 2; ; i++ {
					route.Name = fmt.Sprintf("%s %s %s (%d)", gs.Date, r.Color, r.Grade, i)
					if _, exists := routes[route.Name]; !exists {
						break
					}
				}
			}
			if d.Exists(route.Name, "routes") {
				bug.UserBugOn(len(routes) > 0, fmt.Sprintf("Existing route '%s' from file '%s' detected in database", route.Name, path))
				fmt.Printf("Existing route '%s' from file '%s' in the database, skipping entire file.\n", route.Name, path)
				return
			}
			route.Grade, route.Type = getGymGrade(r.Grade)
			route.Stars = getStars(route.Name, r.Stars)

			tick := &database.Tick{
				CragId:   gym.Id,
				AreaId:   wall.Id,
				Date:     date,
				Attempts: default1(r.Attempts),
				Sessions: default1(r.Sessions),
				Lead:     getGymLead(r.TopRope, route.Type),
				Falls:    r.Falls,
				Hangs:    r.Hangs,
			}
			tick.Redpoint, tick.Flash, tick.Onsight = getRFO(tick.Attempts, r.Falls, r.Hangs, r.Flash, r.Onsight)
			routes[route.Name] = RouteAndTick{route, tick}
		}
	}

	bug.UserBugOn(len(routes) == 0, fmt.Sprintf("no routes defined in '%s'", path))

	lps := make([]*database.DoubleLP, 0, len(routes))
	for _, x := range routes {
		lps = append(lps, &database.DoubleLP{Side1: x.route, Side2: x.tick})
	}
	d.InsertDoubleLPs(lps)
}
