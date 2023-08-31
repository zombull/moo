package database

import (
	"database/sql"

	"github.com/zombull/moo/bug"
)

type Setter struct {
	Id       int64  `yaml:"-"`
	Name     string `yaml:"name"`
	Nickname string `yaml:"nickname"`
}

const SETTER_SCHEMA string = `
CREATE TABLE IF NOT EXISTS setters (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	nickname TEXT,
	UNIQUE (name)
);
`

func (s *Setter) id() int64 {
	return s.Id
}

func (s *Setter) idName() string {
	return "id"
}

func (s *Setter) setId(id int64) {
	s.Id = id
}

func (s *Setter) setSideOneId(id int64) {
	bug.Bug("Double LP is only for Problem+Holds")
}

func (s *Setter) table() string {
	return "setters"
}

func (s *Setter) keys() []string {
	return []string{"name", "nickname"}
}

func (s *Setter) values() []interface{} {
	return []interface{}{s.Name, s.Nickname}
}

func (d *Database) scanSetters(r *sql.Rows) []*Setter {
	defer r.Close()

	var setters []*Setter
	for r.Next() {
		s := Setter{}
		err := r.Scan(
			&s.Id,
			&s.Name,
			&s.Nickname,
		)
		bug.OnError(err)
		setters = append(setters, &s)
	}

	return setters
}

func (d *Database) FindSetter(name string) *Setter {
	r := d.query(`SELECT * FROM setters WHERE name=?`, []interface{}{name})
	s := d.scanSetters(r)
	if len(s) == 0 {
		return nil
	}
	return s[0]
}

func (d *Database) GetSetter(id int64) *Setter {
	r := d.query(`SELECT * FROM setters WHERE id=?`, []interface{}{id})
	s := d.scanSetters(r)
	bug.On(len(s) == 0, sql.ErrNoRows.Error())
	return s[0]
}

func (d *Database) GetSetters() []*Setter {
	r := d.query(`SELECT * FROM setters`, []interface{}{})
	return d.scanSetters(r)
}
