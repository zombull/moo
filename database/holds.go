package database

import (
	"database/sql"

	"github.com/zombull/moo/bug"
)

type Holds struct {
	ProblemId int64  `yaml:"-"`
	Holds   []string `yaml:"holds"`
}

var HoldKeys = []string{
	"problem_id",
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
	"h7",
	"h8",
	"h9",
	"h10",
	"h11",
	"h12",
	"h13",
	"h14",
	"h15",
	"h16",
	"h17",
	"h18",
	"h19",
	"h20",
	"h21",
	"h22",
	"h23",
	"h24",
}

const HOLDS_SCHEMA = `
CREATE TABLE IF NOT EXISTS holds (
	problem_id INTEGER PRIMARY KEY,
	h1 TEXT NOT NULL,
	h2 TEXT NOT NULL,
	h3 TEXT,
	h4 TEXT,
	h5 TEXT,
	h6 TEXT,
	h7 TEXT,
	h8 TEXT,
	h9 TEXT,
	h10 TEXT,
	h11 TEXT,
	h12 TEXT,
	h13 TEXT,
	h14 TEXT,
	h15 TEXT,
	h16 TEXT,
	h17 TEXT,
	h18 TEXT,
	h19 TEXT,
	h20 TEXT,
	h21 TEXT,
	h22 TEXT,
	h23 TEXT,
	h24 TEXT,
	FOREIGN KEY (problem_id) REFERENCES problems (id)
);`

func (h *Holds) id() int64 {
	return h.ProblemId
}

func (h *Holds) idName() string {
	return "problem_id"
}

func (h *Holds) setSideOneId(id int64) {
	h.ProblemId = id
}

func (h *Holds) setId(id int64) {
	h.ProblemId = id
}

func (h *Holds) table() string {
	return "holds"
}

func (h *Holds) keys() []string {
	return HoldKeys
}

func (h *Holds) values() []interface{} {
	values := make([]interface{}, len(HoldKeys))
	for i := range HoldKeys {
		if i == 0 {
			values[i] = h.ProblemId
		} else if i-1 < len(h.Holds) {
			values[i] = h.Holds[i-1]
		}
	}
	return values
}

func (d *Database) scanHolds(r *sql.Rows) *Holds {
	defer r.Close()

	for r.Next() {
		// values := make([]interface{}, len(HoldKeys))

		id := int64(-1)
		s := make([]sql.NullString, len(HoldKeys)-1)
		err := r.Scan(
			&id,
			&s[0],
			&s[1],
			&s[2],
			&s[3],
			&s[4],
			&s[5],
			&s[6],
			&s[7],
			&s[8],
			&s[9],
			&s[10],
			&s[11],
			&s[12],
			&s[13],
			&s[14],
			&s[15],
			&s[16],
			&s[17],
			&s[18],
			&s[19],
			&s[20],
			&s[21],
			&s[22],
			&s[23],
		)
		bug.OnError(err)

		h := &Holds{ProblemId: id, Holds: make([]string, 0, 2)}
		for i := range HoldKeys {
			if i != 0 && s[i-1].Valid {
				h.Holds = append(h.Holds, s[i-1].String)
			}
		}
		return h
	}
	return nil
}

func (d *Database) GetHolds(id int64) *Holds {
	q := "SELECT * FROM holds WHERE problem_id=?"
	h := d.query(q, []interface{}{id})
	holds := d.scanHolds(h)
	bug.On(holds == nil, sql.ErrNoRows.Error())
	return holds
}
