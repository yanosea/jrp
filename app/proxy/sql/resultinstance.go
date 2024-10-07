package sqlproxy

import (
	"database/sql"
)

// ResultInstanceInterface is an interface for sql.Result.
type ResultInstanceInterface interface {
	RowsAffected() (int64, error)
}

// ResultInstance is a struct that implements ResultInstanceInterface.
type ResultInstance struct {
	FieldResult sql.Result
}

// RowsAffected is a proxy for sql.Result.RowsAffected.
func (r ResultInstance) RowsAffected() (int64, error) {
	return r.FieldResult.RowsAffected()
}
