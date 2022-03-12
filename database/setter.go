package database

import (
	"database/sql"
	"fmt"

	"github.com/zombull/moo/bug"
)

type Setter struct {
	Id       int64  `yaml:"-"`
	CragId   int64  `yaml:"-"`
	Name     string `yaml:"name"`
	Nickname string `yaml:"nickname"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Inactive bool   `yaml:"inactive"`
	Comment  string `yaml:"name"`
}

const SETTER_SCHEMA string = `
CREATE TABLE IF NOT EXISTS setters (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	crag_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	nickname TEXT,
	country TEXT NOT NULL,
	city TEXT NOT NULL,
	inactive BOOLEAN NOT NULL,
	comment TEXT,
	FOREIGN KEY (crag_id) REFERENCES crags (id),
	UNIQUE (crag_id, name)
);
`

func (s *Setter) id() int64 {
	return s.Id
}

func (s *Setter) setId(id int64) {
	s.Id = id
}

func (s *Setter) setSideOneId(id int64) {
	s.CragId = id
}

func (s *Setter) table() string {
	return "setters"
}

func (s *Setter) keys() []string {
	return []string{"crag_id", "name", "nickname", "country", "city", "inactive", "comment"}
}

func (s *Setter) values() []interface{} {
	return []interface{}{s.CragId, s.Name, s.Nickname, s.Country, s.City, s.Inactive, s.Comment}
}

func (d *Database) DeleteSetter(id int64) {
	panic(fmt.Errorf("DeleteSetter not yet implemented"))
}

func (d *Database) scanSetters(r *sql.Rows) []*Setter {
	defer r.Close()

	var setters []*Setter
	for r.Next() {
		s := Setter{}
		err := r.Scan(
			&s.Id,
			&s.CragId,
			&s.Name,
			&s.Nickname,
			&s.Country,
			&s.City,
			&s.Inactive,
			&s.Comment,
		)
		bug.OnError(err)
		setters = append(setters, &s)
	}

	return setters
}

func (d *Database) FindSetter(cragId int64, name string) *Setter {
	r := d.query(`SELECT * FROM setters WHERE crag_id=? AND name=?`, []interface{}{cragId, name})
	a := d.scanSetters(r)
	if len(a) == 0 {
		return nil
	}
	return a[0]
}

func (d *Database) GetSetter(id int64) *Setter {
	r := d.query(`SELECT * FROM setters WHERE id=?`, []interface{}{id})
	a := d.scanSetters(r)
	bug.On(len(a) == 0, sql.ErrNoRows.Error())
	return a[0]
}

func (d *Database) GetSetters(cragId int64) []*Setter {
	r := d.query(`SELECT * FROM setters WHERE crag_id=?`, []interface{}{cragId})
	return d.scanSetters(r)
}
