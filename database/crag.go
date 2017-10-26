package database

import (
	"database/sql"

	"github.com/zombull/floating-castle/bug"
)

type Crag struct {
	Id       int64  `yaml:"-"`
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	Url      string `yaml:"url"`
	Map      string `yaml:"map"`
	Comment  string `yaml:"comment"`
}

const CRAG_SCHEMA string = `
CREATE TABLE IF NOT EXISTS crags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	location TEXT,
	url TEXT,
	map TEXT,
	comment TEXT,
	UNIQUE (name)
);
`

func (c *Crag) id() int64 {
	return c.Id
}

func (c *Crag) setId(id int64) {
	c.Id = id
}

func (c *Crag) table() string {
	return "crags"
}

func (c *Crag) keys() []string {
	return []string{"name", "location", "url", "map", "comment"}
}

func (c *Crag) values() []interface{} {
	return []interface{}{c.Name, c.Location, c.Url, c.Map, c.Comment}
}

func (d *Database) GetCragId(name string) int64 {
	return d.GetId("crags", name)
}

func (d *Database) scanCrags(r *sql.Rows) []*Crag {
	defer r.Close()

	var crags []*Crag
	for r.Next() {
		c := Crag{}
		err := r.Scan(
			&c.Id,
			&c.Name,
			&c.Location,
			&c.Url,
			&c.Map,
			&c.Comment,
		)
		bug.OnError(err)
		crags = append(crags, &c)
	}

	return crags
}

func (d *Database) FindCrag(name string) *Crag {
	r := d.query("SELECT * FROM crags WHERE name=?", []interface{}{name})
	c := d.scanCrags(r)
	if len(c) == 0 {
		return nil
	}
	return c[0]
}

func (d *Database) GetCrag(id int64) *Crag {
	r := d.query(`SELECT * FROM crags WHERE id=?`, []interface{}{id})
	c := d.scanCrags(r)
	bug.On(len(c) == 0, sql.ErrNoRows.Error())
	return c[0]
}

func (d *Database) GetCrags() []*Crag {
	r := d.query(`SELECT * FROM crags`, []interface{}{})
	return d.scanCrags(r)
}
