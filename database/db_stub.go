// +build server

package database

import (
	"database/sql"

	"github.com/zombull/moo/bug"
)

type Database struct {
	db *sql.DB
}

type Record interface {
	id() int64
	idName() string
	setId(int64)
	table() string
	keys() []string
	values() []interface{}
}

type SideTwo interface {
	Record
	setSideOneId(int64)
}

type DoubleLP struct {
	Side1 Record
	Side2 SideTwo
}

func Init(path string) *Database {
	bug.Bug("Attempting to instantiate stubbed DB\n");
	return nil
}

func (d *Database) Insert(r Record) {
	bug.Bug("Attempting to insert into stubbed DB\n");
}

func (d *Database) InsertDoubleLP(side1 Record, side2 SideTwo) {
	bug.Bug("Attempting to insert into stubbed DB\n");
}

func (d *Database) InsertDoubleLPs(lps []*DoubleLP) {
	bug.Bug("Attempting to insert into stubbed DB\n");
}

func (d *Database) InsertCollection(master Record, records []SideTwo) {
	bug.Bug("Attempting to insert into stubbed DB\n");
}

func (d *Database) Update(r Record) {
	bug.Bug("Attempting to update into stubbed DB\n");
}

func (d *Database) Delete(r Record) {
	bug.Bug("Attempting to delete from stubbed DB\n");
}

func (d *Database) GetId(table, name string) int64 {
	bug.Bug("Attempting to read from stubbed DB\n");
	return int64(-1)
}

func (d *Database) queryRow(q string, args []interface{}) *sql.Row {
	bug.Bug("Attempting to read from stubbed DB\n");
	return nil
}

func (d *Database) query(q string, args []interface{}) *sql.Rows {
	bug.Bug("Attempting to read from stubbed DB\n");
	return nil
}

func (d *Database) ExistsBy(field, value, table string) bool {
	bug.Bug("Attempting to read from stubbed DB\n");
	return false
}

func (d *Database) Exists(name, table string) bool {
	return d.ExistsBy("name", name, table)
}
