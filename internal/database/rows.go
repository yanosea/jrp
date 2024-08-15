package database

import (
	"database/sql"
)

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

type SQLiteRows struct {
	*sql.Rows
}

func (s SQLiteRows) Next() bool {
	return s.Rows.Next()
}

func (s SQLiteRows) Scan(dest ...any) error {
	return s.Rows.Scan(dest...)
}

func (s SQLiteRows) Close() error {
	return s.Rows.Close()
}
