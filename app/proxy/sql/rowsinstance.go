package sqlproxy

import (
	"database/sql"
)

// RowsInstanceInterface is an interface for sql.Rows.
type RowsInstanceInterface interface {
	Close() error
	Next() bool
	Scan(dest ...any) error
}

// RowsInstance is a struct that implements RowsInstanceInterface.
type RowsInstance struct {
	FieldRows *sql.Rows
}

// Close is a proxy for sql.Rows.Close.
func (r *RowsInstance) Close() error {
	return r.FieldRows.Close()
}

// Next is a proxy for sql.Rows.Next.
func (r *RowsInstance) Next() bool {
	return r.FieldRows.Next()
}

// Scan is a proxy for sql.Rows.Scan.
func (r *RowsInstance) Scan(dest ...any) error {
	return r.FieldRows.Scan(dest...)
}
