package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type DatabaseProvider interface {
	Connect(dbFilePath string) (*sql.DB, error)
	Query(db *sql.DB, query string) (*sql.Rows, error)
}

type SQLiteProvider struct{}

func (s SQLiteProvider) Connect(dbFilePath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbFilePath)
}

func (s SQLiteProvider) Query(db *sql.DB, query string) (*sql.Rows, error) {
	return db.Query(query)
}
