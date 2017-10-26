package interactive

import (
	"github.com/zombull/floating-castle/database"
)

type queryOp struct {
}

func (q *queryOp) name() string {
	return "Query"
}

func (q *queryOp) crag(d *database.Database, p string) {

}

func (q *queryOp) area(d *database.Database, p string) {

}

func (q *queryOp) setter(d *database.Database) {

}

func (q *queryOp) route(d *database.Database) {

}

func (q *queryOp) tick(d *database.Database) {

}

func (q *queryOp) list(d *database.Database) {

}
