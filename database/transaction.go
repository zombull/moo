package database

import (
	"database/sql"

	"github.com/zombull/moo/bug"
)

type transaction struct {
	tx *sql.Tx
}

func (t *transaction) prepare(query string) *statement {
	stmt, err := t.tx.Prepare(query)
	if err != nil {
		t.tx.Rollback()
		bug.OnError(err)
	}
	return &statement{stmt}
}

func (t *transaction) exec(query string, args ...interface{}) sql.Result {
	r, err := t.tx.Exec(query, args...)
	if err != nil {
		t.tx.Rollback()
		bug.OnError(err)
	}
	return r
}

func (t *transaction) getInsertedId(r sql.Result) int64 {
	id, err := r.LastInsertId()
	if err != nil {
		t.tx.Rollback()
		bug.OnError(err)
	}
	return id
}
