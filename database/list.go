package database

import (
	"fmt"

	"github.com/zombull/floating-castle/bug"
)

type List struct {
	Name   string
	Routes []*Route
}

const LIST_SCHEMA string = `
CREATE TABLE IF NOT EXISTS lists (
	name TEXT PRIMARY KEY
);
`

func (d *Database) InsertList(name string) {
	bug.On(d.Exists(name, "lists"), fmt.Sprintf("list '%s' already exists", name))

	tx := d.begin()

	stmt := tx.prepare(`INSERT INTO lists (name) VALUES (?)`)
	defer stmt.close()

	stmt.exec(tx, name)

	cstmt := tx.prepare(fmt.Sprintf("CREATE TABLE list_%s (route_id INTEGER PRIMARY KEY);", name))
	defer cstmt.close()

	cstmt.exec(tx)
	d.commit(tx)
}

func (d *Database) DeleteList(id int) {
	panic(fmt.Errorf("DeleteCrag not yet implemented"))
}

func (d *Database) GetLists() []string {
	q := "SELECT name FROM lists"
	r := d.query(q, []interface{}{})
	defer r.Close()

	var names []string
	for r.Next() {
		name := ""
		err := r.Scan(&name)
		bug.OnError(err)
		names = append(names, name)
	}
	return names
}

func (d *Database) GetList(name string) *List {
	if !d.Exists(name, "lists") {
		panic(fmt.Errorf("The list '%s' does not exist", name))
	}

	q := fmt.Sprintf("SELECT id FROM list_%s", name)
	r := d.query(q, []interface{}{})

	l := List{Name: name}
	for r.Next() {
		id := int64(-1)
		if err := r.Scan(&id); err != nil {
			panic(err)
		}
		l.Routes = append(l.Routes, d.GetRoute(id))
	}
	return &l
}
