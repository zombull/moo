package database

import (
	"database/sql"

	"github.com/zombull/moo/bug"
)

var Stars = map[string]int{
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
}

type Problem struct {
	Id        int64     `yaml:"-"`
	SetId     int64     `yaml:"-"`
	SetterId  int64     `yaml:"-"`
	Name      string    `yaml:"name"`
	Type      string    `yaml:"type"`
	Date      int64     `yaml:"date"`
	Grade     string    `yaml:"grade"`
	Stars     uint      `yaml:"stars"`
	MoonId    uint      `yaml:"moon_id"`
	Ascents   uint      `yaml:"ascents"`
	Benchmark bool      `yaml:"benchmark"`
}

const PROBLEM_SCHEMA string = `
CREATE TABLE IF NOT EXISTS problems (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	set_id INTEGER NOT NULL,
	setter_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	date NUMERIC NOT NULL,
	grade TEXT NOT NULL,
	stars INTEGER NOT NULL,
	moon_id INTEGER NOT NULL,
	ascents INTEGER NOT NULL,
	benchmark BOOLEAN,
	FOREIGN KEY (set_id) REFERENCES sets (id),
	FOREIGN KEY (setter_id) REFERENCES setters (id),
	UNIQUE (set_id, name)
);
`

func (p *Problem) id() int64 {
	return p.Id
}

func (p *Problem) idName() string {
	return "id"
}

func (p *Problem) setId(id int64) {
	p.Id = id
}

func (p *Problem) table() string {
	return "problems"
}

func (p *Problem) keys() []string {
	return []string{"set_id", "setter_id", "name", "date", "grade", "stars", "moon_id", "ascents", "benchmark"}
}

func (p *Problem) values() []interface{} {
	return []interface{}{p.SetId, p.SetterId, p.Name, p.Date, p.Grade, p.Stars, p.MoonId, p.Ascents, p.Benchmark}
}

func (d *Database) scanProblems(r *sql.Rows) []*Problem {
	defer r.Close()

	var problems []*Problem
	for r.Next() {
		problem := Problem{}
		err := r.Scan(
			&problem.Id,
			&problem.SetId,
			&problem.SetterId,
			&problem.Name,
			&problem.Date,
			&problem.Grade,
			&problem.Stars,
			&problem.MoonId,
			&problem.Ascents,
			&problem.Benchmark,
		)
		bug.OnError(err)
		problems = append(problems, &problem)
	}

	return problems
}

func (d *Database) GetProblems(setId int64) []*Problem {
	r := d.query(`SELECT * FROM problems WHERE set_id=?`, []interface{}{setId})
	return d.scanProblems(r)
}

func (d *Database) GetProblem(id int64) *Problem {
	q := "SELECT * FROM problems WHERE id=?"
	r := d.query(q, []interface{}{id})
	problems := d.scanProblems(r)
	bug.On(len(problems) == 0, sql.ErrNoRows.Error())
	return problems[0]
}

func (d *Database) FindProblem(setId int64, name string) *Problem {
	q := `SELECT * FROM problems WHERE set_id=? AND name=?`
	r := d.query(q, []interface{}{setId, name})
	problems := d.scanProblems(r)
	if len(problems) == 0 {
		return nil
	}
	return problems[0]
}

func (d *Database) FindProblemByMoonId(setId int64, moonId uint) *Problem {
	q := `SELECT * FROM problems WHERE set_id=? AND moon_id=?`
	r := d.query(q, []interface{}{setId, moonId})
	problems := d.scanProblems(r)
	if len(problems) == 0 {
		return nil
	}
	return problems[0]
}
