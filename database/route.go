package database

import (
	"database/sql"
	"time"

	"github.com/zombull/moo/bug"
)

var RouteTypes = Set{
	"moonboard": xx,
	"boulder":   xx,
	"sport":     xx,
	"trad":      xx,
}

var Stars = map[string]int{
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
}

type Route struct {
	Id        int64     `yaml:"-"`
	CragId    int64     `yaml:"-"`
	AreaId    int64     `yaml:"-"`
	SetterId  int64     `yaml:"-"`
	Name      string    `yaml:"name"`
	Type      string    `yaml:"type"`
	Date      time.Time `yaml:"date"`
	Grade     string    `yaml:"grade"`
	Stars     uint      `yaml:"stars"`
	Length    uint      `yaml:"length"`  // doubles as Moonboard ID
	Pitches   uint      `yaml:"pitches"` // doubles as Moonboard ascents/repeats
	Benchmark bool      `yaml:"benchmark"`
	Url       string    `yaml:"url"`
	Comment   string    `yaml:"comment"`
}

const FORMAT_ROUTE = `:
    Name:       %s
    Type:       %s
    Date:       %s
	Grade:      %s
	Benchmark:  %t
    Stars:      %d
    Length:     %d
    Pitches:    %d
    Setter:     %s
    Url:        %s
    Comment:    %s
`

const ROUTE_SCHEMA string = `
CREATE TABLE IF NOT EXISTS routes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	crag_id INTEGER NOT NULL,
	area_id INTEGER NOT NULL,
	setter_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	date DATE NOT NULL,
	grade TEXT NOT NULL,
	stars INTEGER NOT NULL,
	length INTEGER NOT NULL,
	pitches INTEGER NOT NULL,
	benchmark BOOLEAN,
	url TEXT,
	comment TEXT,
	FOREIGN KEY (crag_id) REFERENCES crags (id),
	FOREIGN KEY (area_id) REFERENCES areas (id),
	FOREIGN KEY (setter_id) REFERENCES setters (id),
	UNIQUE (crag_id, area_id, name)
);
`

func (r *Route) id() int64 {
	return r.Id
}

func (r *Route) setId(id int64) {
	r.Id = id
}

func (r *Route) table() string {
	return "routes"
}

func (r *Route) keys() []string {
	return []string{"crag_id", "area_id", "setter_id", "name", "type", "date", "grade", "stars", "length", "pitches", "benchmark", "url", "comment"}
}

func (r *Route) values() []interface{} {
	return []interface{}{r.CragId, r.AreaId, r.SetterId, r.Name, r.Type, r.Date.Unix(), r.Grade, r.Stars, r.Length, r.Pitches, r.Benchmark, r.Url, r.Comment}
}

func (d *Database) scanRoutes(r *sql.Rows) []*Route {
	defer r.Close()

	var routes []*Route
	for r.Next() {
		route := Route{}
		err := r.Scan(
			&route.Id,
			&route.CragId,
			&route.AreaId,
			&route.SetterId,
			&route.Name,
			&route.Type,
			&route.Date,
			&route.Grade,
			&route.Stars,
			&route.Length,
			&route.Pitches,
			&route.Benchmark,
			&route.Url,
			&route.Comment,
		)
		bug.OnError(err)
		routes = append(routes, &route)
	}

	return routes
}

func (d *Database) GetAllRoutes(cragId int64) []*Route {
	r := d.query(`SELECT * FROM routes WHERE crag_id=?`, []interface{}{cragId})
	return d.scanRoutes(r)
}

func (d *Database) GetRoutes(cragId, areaId int64) []*Route {
	r := d.query(`SELECT * FROM routes WHERE crag_id=? AND area_id=?`, []interface{}{cragId, areaId})
	return d.scanRoutes(r)
}

func (d *Database) GetRoute(id int64) *Route {
	q := "SELECT * FROM routes WHERE id=?"
	r := d.query(q, []interface{}{id})
	routes := d.scanRoutes(r)
	bug.On(len(routes) == 0, sql.ErrNoRows.Error())
	return routes[0]
}

func (d *Database) FindRoute(areaId int64, name string) *Route {
	q := `SELECT * FROM routes WHERE area_id=? AND name=?`
	r := d.query(q, []interface{}{areaId, name})
	routes := d.scanRoutes(r)
	if len(routes) == 0 {
		return nil
	}
	return routes[0]
}

func (d *Database) FindRouteByLength(areaId int64, length uint) *Route {
	q := `SELECT * FROM routes WHERE area_id=? AND length=?`
	r := d.query(q, []interface{}{areaId, length})
	routes := d.scanRoutes(r)
	if len(routes) == 0 {
		return nil
	}
	return routes[0]
}
