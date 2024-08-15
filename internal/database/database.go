package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type DatabaseProvider interface {
	Connect(dbFilePath string) (*sql.DB, error)
	Query(db *sql.DB, query string) (Rows, error)
}

type SQLiteProvider struct{}

func (s SQLiteProvider) Connect(dbFilePath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbFilePath)
}

func (s SQLiteProvider) Query(db *sql.DB, query string) (Rows, error) {
	rows, _ := db.Query(query)
	return &SQLiteRows{rows}, nil
}
