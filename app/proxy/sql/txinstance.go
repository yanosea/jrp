package sqlproxy

import (
	"database/sql"
)

// TxInstanceInterface is an interface for sql.Tx.
type TxInstanceInterface interface {
	Commit() error
	Exec(query string, args ...interface{}) (ResultInstanceInterface, error)
	QueryRow(query string, args ...interface{}) RowInstanceInterface
	Rollback() error
}

// TxInstance is a struct that implements TxInstanceInterface.
type TxInstance struct {
	FieldTx *sql.Tx
}

// Commit is a proxy for sql.Tx.Commit.
func (t *TxInstance) Commit() error {
	return t.FieldTx.Commit()
}

// Exec is a proxy for sql.Tx.Exec.
func (t *TxInstance) Exec(query string, args ...interface{}) (ResultInstanceInterface, error) {
	res, err := t.FieldTx.Exec(query, args...)
	return &ResultInstance{FieldResult: res}, err
}

// QueryRow is a proxy for sql.Tx.QueryRow.
func (t *TxInstance) QueryRow(query string, args ...interface{}) RowInstanceInterface {
	return &RowInstance{FieldRow: t.FieldTx.QueryRow(query, args...)}
}

// Rollback is a proxy for sql.Tx.Rollback.
func (t *TxInstance) Rollback() error {
	return t.FieldTx.Rollback()
}
