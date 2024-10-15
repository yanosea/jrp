package sqlproxy

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Sql is an interface for sql.
type Sql interface {
	Open(driverName string, dataSourceName string) (DBInstanceInterface, error)
	StringToNullString(s string) *NullStringInstance
}

// SqlProxy is a struct that implements Sql.
type SqlProxy struct{}

// New is a constructor for SqlProxy.
func New() Sql {
	return &SqlProxy{}
}

// Open is a proxy for sql.Open.
func (*SqlProxy) Open(driverName string, dataSourceName string) (DBInstanceInterface, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DBInstance{FieldDB: db}, err
}

// StringToNullString returns a NullStringInstance with the argument as the String field.
func (*SqlProxy) StringToNullString(s string) *NullStringInstance {
	return &NullStringInstance{
		FieldNullString: &sql.NullString{
			String: s,
			Valid:  s != "",
		},
	}
}
