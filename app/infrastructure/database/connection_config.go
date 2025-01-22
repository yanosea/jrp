package database

import ()

// DBType is a type that represents the type of database.
type DBType string

const (
	// SQLite is a type of database.
	SQLite DBType = "sqlite"
)

// DBName is a type that represents the name of the database.
type DBName string

const (
	// JrpDB is the name of the Jrp database.
	JrpDB DBName = "jrp"
	// WNJpnDB is the name of the WNJpn database.
	WNJpnDB DBName = "wnjpn"
)

// ConnectionConfig is a struct that contains the configuration of the database connection.
type ConnectionConfig struct {
	DBName DBName
	DBType DBType
	DSN    string
}
