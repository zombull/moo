package database

import (
	"database/sql"
	"time"

	"github.com/zombull/moo/bug"
)

type Tick struct {
	Id       int64     `yaml:"-"`
	RouteId  int64     `yaml:"-"`
	AreaId   int64     `yaml:"-"`
	CragId   int64     `yaml:"-"`
	Date     time.Time `yaml:"date"`
	Grade    string    `yaml:"grade"`
	Stars    uint      `yaml:"stars"`
	Attempts uint      `yaml:"attempts"`
	Sessions uint      `yaml:"sessions"`
	Redpoint bool      `yaml:"redpoint"`
	Flash    bool      `yaml:"flash"`
	Onsight  bool      `yaml:"onsight"`
	Lead     bool      `yaml:"lead"`
	Falls    uint      `yaml:"false"`
	Hangs    uint      `yaml:"hangs"`
	Comment  string    `yaml:"comment"`
}

const FORMAT_TICK = `:
    Name:     %s
	Date:     %s
	Grade:    %s
	Stars:    %d
    Redpoint: %t
    Flash:    %t
    Onsight:  %t
    Falls:    %d
    Hangs:    %d
    Attempts: %d
    Sessions: %d
    Comment:  %s
`

const TICK_SCHEMA string = `
CREATE TABLE IF NOT EXISTS ticks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	crag_id INTEGER NOT NULL,
	area_id INTEGER NOT NULL,
	route_id INTEGER NOT NULL,
	date DATE NOT NULL,
	grade TEXT NOT NULL,
	stars INTEGER NOT NULL,
	lead BOOLEAN NOT NULL,
	redpoint BOOLEAN NOT NULL,
	flash BOOLEAN NOT NULL,
	onsight BOOLEAN NOT NULL,
	falls INTEGER NOT NULL,
	hangs INTEGER NOT NULL,
	attempts INTEGER NOT NULL,
	sessions INTEGER NOT NULL,
	comment TEXT,
	FOREIGN KEY (route_id) REFERENCES routes (id),
	FOREIGN KEY (area_id) REFERENCES areas (id),
	FOREIGN KEY (crag_id) REFERENCES crags (id)
);
`

func (t *Tick) id() int64 {
	return t.Id
}

func (t *Tick) setSideOneId(id int64) {
	t.RouteId = id
}

func (t *Tick) setId(id int64) {
	t.Id = id
}

func (t *Tick) table() string {
	return "ticks"
}

func (t *Tick) keys() []string {
	return []string{"crag_id", "area_id", "route_id", "date", "grade", "stars", "lead", "redpoint", "flash", "onsight", "falls", "hangs", "attempts", "sessions", "comment"}
}

func (t *Tick) values() []interface{} {
	return []interface{}{t.CragId, t.AreaId, t.RouteId, t.Date.Unix(), t.Grade, t.Stars, t.Lead, t.Redpoint, t.Flash, t.Onsight, t.Falls, t.Hangs, t.Attempts, t.Sessions, t.Comment}
}

func (d *Database) scanTicks(r *sql.Rows) []*Tick {
	defer r.Close()

	var ticks []*Tick
	for r.Next() {
		t := Tick{}
		err := r.Scan(
			&t.Id,
			&t.CragId,
			&t.AreaId,
			&t.RouteId,
			&t.Date,
			&t.Grade,
			&t.Stars,
			&t.Lead,
			&t.Redpoint,
			&t.Flash,
			&t.Onsight,
			&t.Falls,
			&t.Hangs,
			&t.Attempts,
			&t.Sessions,
			&t.Comment,
		)
		bug.OnError(err)
		ticks = append(ticks, &t)
	}

	return ticks
}

func (d *Database) GetTicks(routeId int64) []*Tick {
	r := d.query(`SELECT * FROM ticks WHERE route_id=?`, []interface{}{routeId})
	return d.scanTicks(r)
}

func (d *Database) GetAreaTicks(areaId int64) []*Tick {
	r := d.query(`SELECT * FROM ticks WHERE area_id=?`, []interface{}{areaId})
	return d.scanTicks(r)
}

func (d *Database) GetCragTicks(cragId int64) []*Tick {
	r := d.query(`SELECT * FROM ticks WHERE crag_id=?`, []interface{}{cragId})
	return d.scanTicks(r)
}

func (d *Database) GetAllTicks() []*Tick {
	r := d.query(`SELECT * FROM ticks`, []interface{}{})
	return d.scanTicks(r)
}
