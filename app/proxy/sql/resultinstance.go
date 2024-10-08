package sqlproxy

import (
	"database/sql"
)

// ResultInstanceInterface is an interface for sql.Result.
type ResultInstanceInterface interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// ResultInstance is a struct that implements ResultInstanceInterface.
type ResultInstance struct {
	FieldResult sql.Result
}

// LastInsertId is a proxy for sql.Result.LastInsertId.
func (r ResultInstance) LastInsertId() (int64, error) {
	return r.FieldResult.LastInsertId()
}

// RowsAffected is a proxy for sql.Result.RowsAffected.
func (r ResultInstance) RowsAffected() (int64, error) {
	return r.FieldResult.RowsAffected()
}
