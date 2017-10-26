package database

import (
	"database/sql"

	"github.com/zombull/floating-castle/bug"
)

type Area struct {
	Id      int64  `yaml:"-"`
	CragId  int64  `yaml:"-"`
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	Map     string `yaml:"map"`
	Comment string `yaml:"comment"`
}

const AREA_SCHEMA string = `
CREATE TABLE IF NOT EXISTS areas (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	crag_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	url TEXT,
	map TEXT,
	comment TEXT,	
	FOREIGN KEY (crag_id) REFERENCES crags (id),
	UNIQUE (crag_id, name)
);
`

func (a *Area) id() int64 {
	return a.Id
}

func (a *Area) setId(id int64) {
	a.Id = id
}

func (a *Area) setSideOneId(id int64) {
	a.CragId = id
}

func (a *Area) table() string {
	return "areas"
}

func (a *Area) keys() []string {
	return []string{"crag_id", "name", "url", "map", "comment"}
}

func (a *Area) values() []interface{} {
	return []interface{}{a.CragId, a.Name, a.Url, a.Map, a.Comment}
}

func (d *Database) GetAreaId(cragId int64, name string) int64 {
	id := int64(-1)
	q := `SELECT id FROM areas WHERE name=? AND crag_id=?`
	err := d.queryRow(q, []interface{}{name, cragId}).Scan(&id)
	bug.OnError(err)
	return id
}

func (d *Database) scanAreas(r *sql.Rows) []*Area {
	defer r.Close()

	var areas []*Area
	for r.Next() {
		a := Area{}
		err := r.Scan(
			&a.Id,
			&a.CragId,
			&a.Name,
			&a.Url,
			&a.Map,
			&a.Comment,
		)
		bug.OnError(err)
		areas = append(areas, &a)
	}

	return areas
}

func (d *Database) FindArea(cragId int64, name string) *Area {
	r := d.query(`SELECT * FROM areas WHERE crag_id=? AND name=?`, []interface{}{cragId, name})
	a := d.scanAreas(r)
	if len(a) == 0 {
		return nil
	}
	return a[0]
}

func (d *Database) GetArea(id int64) *Area {
	r := d.query(`SELECT * FROM areas WHERE id=?`, []interface{}{id})
	a := d.scanAreas(r)
	bug.On(len(a) == 0, sql.ErrNoRows.Error())
	return a[0]
}

func (d *Database) GetAreas(cragId int64) []*Area {
	r := d.query(`SELECT * FROM areas WHERE crag_id=?`, []interface{}{cragId})
	return d.scanAreas(r)
}
