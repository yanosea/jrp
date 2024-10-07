package sqlproxy

import (
	"database/sql"
)

// RowInstanceInterface is an interface for sql.Row.
type RowInstanceInterface interface {
	Scan(dest ...interface{}) error
}

// RowInstance is a struct that implements RowInstanceInterface.
type RowInstance struct {
	FieldRow *sql.Row
}

// Scan is a proxy for sql.Row.Scan.
func (r *RowInstance) Scan(dest ...interface{}) error {
	return r.FieldRow.Scan(dest...)
}
