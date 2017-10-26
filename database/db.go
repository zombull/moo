package database

import (
	"database/sql"
	"fmt"
	"path"
	"strings"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/zombull/floating-castle/bug"
)

var TableTypes = Set{
	"crag":   xx,
	"gym":    xx,
	"area":   xx,
	"wall":   xx,
	"setter": xx,
	"route":  xx,
	"tick":   xx,
	"list":   xx,
}

var insertValues = map[int]string{
	1:  "?",
	2:  "?, ?",
	3:  "?, ?, ?",
	4:  "?, ?, ?, ?",
	5:  "?, ?, ?, ?, ?",
	6:  "?, ?, ?, ?, ?, ?",
	7:  "?, ?, ?, ?, ?, ?, ?",
	8:  "?, ?, ?, ?, ?, ?, ?, ?",
	9:  "?, ?, ?, ?, ?, ?, ?, ?, ?",
	10: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	11: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	12: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	13: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	14: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	15: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	16: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	17: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	18: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	19: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	20: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	21: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	22: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	23: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	24: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
	25: "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
}

type Database struct {
	db *sql.DB
}

type Record interface {
	id() int64
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

const SCHEMA string = CRAG_SCHEMA + AREA_SCHEMA + ROUTE_SCHEMA + HOLDS_SCHEMA + TICK_SCHEMA + LIST_SCHEMA + SETTER_SCHEMA

func Init(dir string) *Database {
	timeout := 5 // TODO - make this command-line configurable?

	// These are used to tune the transaction BEGIN behavior instead of using the
	// similar "locking_mode" pragma (locking for the whole database connection).
	path := path.Join(dir, "floating-castle.db")
	path = fmt.Sprintf("%s?_busy_timeout=%d&_txlock=exclusive", path, timeout*1000)

	// Open the database.  Automatically created if it doesn't exist.
	db, err := sql.Open("sqlite3", path)
	bug.On(err != nil, fmt.Sprintf("%v: %s", err, path))

	_, err = db.Exec(SCHEMA)
	bug.OnError(err)

	// Run PRAGMA statements now since they are *per-connection*.
	_, err = db.Exec("PRAGMA foreign_keys=ON;")
	bug.OnError(err)

	return &Database{db}
}

func isLockedError(err error) bool {
	if err == nil {
		return false
	}
	if err == sqlite3.ErrLocked || err == sqlite3.ErrBusy {
		return true
	}
	if err.Error() == "database is locked" {
		return true
	}
	return false
}

func isNoMatchError(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == "sql: no rows in result set" {
		return true
	}
	return false
}

func (d *Database) action(f func() error) error {
	for i := 0; i < 100; i++ {
		err := f()
		if err == nil {
			return nil
		}
		if !isLockedError(err) {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("database deadlock")
}

func (d *Database) begin() *transaction {
	var tx *sql.Tx

	err := d.action(func() error {
		var e2 error
		tx, e2 = d.db.Begin()
		return e2
	})
	bug.OnError(err)
	return &transaction{tx}
}

func (d *Database) commit(t *transaction) {
	err := d.action(t.tx.Commit)
	bug.OnError(err)
}

func (d *Database) transact(query string, getId bool, args ...interface{}) int64 {
	tx := d.begin()

	stmt := tx.prepare(query)
	defer stmt.close()

	r := stmt.exec(tx, args...)

	id := int64(-1)
	if getId {
		id = tx.getInsertedId(r)
	}
	d.commit(tx)
	return id
}

func (d *Database) insert(r Record) string {
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, r.table(), strings.Join(r.keys(), ", "), insertValues[len(r.values())])
}

func (d *Database) Insert(r Record) {
	r.setId(d.transact(d.insert(r), true, r.values()...))
}

func (d *Database) InsertDoubleLP(side1 Record, side2 SideTwo) {
	d.InsertDoubleLPs([]*DoubleLP{&DoubleLP{side1, side2}})
}

func (d *Database) InsertDoubleLPs(lps []*DoubleLP) {
	tx := d.begin()
	for _, lp := range lps {
		stmt := tx.prepare(d.insert(lp.Side1))
		defer stmt.close()

		r := stmt.exec(tx, lp.Side1.values()...)
		lp.Side1.setId(tx.getInsertedId(r))

		lp.Side2.setSideOneId(lp.Side1.id())
		stmt2 := tx.prepare(d.insert(lp.Side2))
		defer stmt2.close()

		r2 := stmt2.exec(tx, lp.Side2.values()...)
		lp.Side2.setId(tx.getInsertedId(r2))
	}
	d.commit(tx)
}

func (d *Database) InsertCollection(master Record, records []SideTwo) {
	tx := d.begin()

	stmt := tx.prepare(d.insert(master))
	defer stmt.close()

	res := stmt.exec(tx, master.values()...)
	master.setId(tx.getInsertedId(res))

	for _, r := range records {
		r.setSideOneId(master.id())
		stmt = tx.prepare(d.insert(r))
		defer stmt.close()

		res = stmt.exec(tx, r.values()...)
		r.setId(tx.getInsertedId(res))
	}
	d.commit(tx)
}

func (d *Database) update(r Record) string {
	return fmt.Sprintf(`UPDATE %s SET %s WHERE id=%d`, r.table(), strings.Join(r.keys(), "=?, ")+"=?", r.id())
}

func (d *Database) Update(r Record) {
	d.transact(d.update(r), false, r.values()...)
}

func (d *Database) Delete(r Record) {
	q := fmt.Sprintf(`DELETE FROM %s WHERE id=%d`, r.table(), r.id())
	d.transact(q, false)
}

func (d *Database) GetId(table, name string) int64 {
	id := int64(-1)
	q := fmt.Sprintf(`SELECT id FROM %s WHERE name=?`, table)
	err := d.queryRow(q, []interface{}{name}).Scan(&id)
	bug.OnError(err)
	return id
}

func (d *Database) queryRow(q string, args []interface{}) *sql.Row {
	var r *sql.Row
	err := d.action(func() error {
		r = d.db.QueryRow(q, args...)
		return nil
	})
	bug.OnError(err)
	return r
}

func (d *Database) query(q string, args []interface{}) *sql.Rows {
	var r *sql.Rows
	err := d.action(func() error {
		var e2 error
		r, e2 = d.db.Query(q, args...)
		return e2
	})
	bug.OnError(err)
	return r
}

func (d *Database) ExistsBy(field, value, table string) bool {
	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?", field, table, field)
	err := d.queryRow(q, []interface{}{value}).Scan(&value)
	if err == sql.ErrNoRows {
		return false
	}
	bug.OnError(err)
	return true
}

func (d *Database) Exists(name, table string) bool {
	return d.ExistsBy("name", name, table)
}
