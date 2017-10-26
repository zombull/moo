package database

import (
	"database/sql"

	"github.com/zombull/floating-castle/bug"
)

type statement struct {
	stmt *sql.Stmt
}

func (s *statement) exec(t *transaction, args ...interface{}) sql.Result {
	r, err := s.stmt.Exec(args...)
	if err != nil {
		t.tx.Rollback()
		bug.OnError(err)
	}
	return r
}

func (s *statement) close() {
	s.stmt.Close()
}
