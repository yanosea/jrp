package sqlproxy

import (
	"database/sql"
)

// StmtInstanceInterface is an interface for sql.Stmt.
type StmtInstanceInterface interface {
	Close() error
	Exec(args ...interface{}) (ResultInstanceInterface, error)
	Query(args ...interface{}) (RowsInstanceInterface, error)
}

// StmtInstance is a struct that implements StmtInstanceInterface.
type StmtInstance struct {
	FieldStmt *sql.Stmt
}

// Close is a proxy for sql.Stmt.Close.
func (s *StmtInstance) Close() error {
	return s.FieldStmt.Close()
}

func (s *StmtInstance) Exec(args ...interface{}) (ResultInstanceInterface, error) {
	res, err := s.FieldStmt.Exec(args...)
	return &ResultInstance{FieldResult: res}, err
}

func (s *StmtInstance) Query(args ...interface{}) (RowsInstanceInterface, error) {
	rows, err := s.FieldStmt.Query(args...)
	return &RowsInstance{FieldRows: rows}, err
}
