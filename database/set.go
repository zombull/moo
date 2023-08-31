package database

import (
	"database/sql"

	"github.com/zombull/moo/bug"
)

type Set struct {
	Id      int64  `yaml:"-"`
	Name    string `yaml:"name"`
}

const SET_SCHEMA string = `
CREATE TABLE IF NOT EXISTS sets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	UNIQUE (name)
);
`

func (s *Set) id() int64 {
	return s.Id
}

func (s *Set) idName() string {
	return "id"
}

func (s *Set) setId(id int64) {
	s.Id = id
}

func (s *Set) setSideOneId(id int64) {
	bug.Bug("Double LP is only for Problem+Holds")
}

func (s *Set) table() string {
	return "sets"
}

func (s *Set) keys() []string {
	return []string{"name"}
}

func (s *Set) values() []interface{} {
	return []interface{}{s.Name}
}

func (d *Database) GetSetId(name string) int64 {
	id := int64(-1)
	q := `SELECT id FROM sets WHERE name=?`
	err := d.queryRow(q, []interface{}{name}).Scan(&id)
	bug.OnError(err)
	return id
}

func (d *Database) scanSets(r *sql.Rows) []*Set {
	defer r.Close()

	var sets []*Set
	for r.Next() {
		s := Set{}
		err := r.Scan(
			&s.Id,
			&s.Name,
		)
		bug.OnError(err)
		sets = append(sets, &s)
	}

	return sets
}

func (d *Database) FindSet(name string) *Set {
	r := d.query(`SELECT * FROM sets WHERE name=?`, []interface{}{name})
	a := d.scanSets(r)
	if len(a) == 0 {
		return nil
	}
	return a[0]
}

func (d *Database) GetSet(id int64) *Set {
	r := d.query(`SELECT * FROM sets WHERE id=?`, []interface{}{id})
	s := d.scanSets(r)
	bug.On(len(s) == 0, sql.ErrNoRows.Error())
	return s[0]
}

func (d *Database) GetSets() []*Set {
	r := d.query(`SELECT * FROM sets`, []interface{}{})
	return d.scanSets(r)
}
