package sqlproxy

import (
	"database/sql"
)

// DBInstanceInterface is an interface for sql.DB.
type DBInstanceInterface interface {
	Begin() (TxInstanceInterface, error)
	Close() error
	Exec(query string, args ...interface{}) (ResultInstanceInterface, error)
	Prepare(query string) (StmtInstanceInterface, error)
	Query(query string, args ...interface{}) (RowsInstanceInterface, error)
}

// DBInstance is a struct that implements DBInstanceInterface.
type DBInstance struct {
	FieldDB *sql.DB
}

// Begin is a proxy for sql.DB.Begin.
func (d *DBInstance) Begin() (TxInstanceInterface, error) {
	tx, err := d.FieldDB.Begin()
	return &TxInstance{FieldTx: tx}, err
}

// Close is a proxy for sql.DB.Close.
func (d *DBInstance) Close() error {
	return d.FieldDB.Close()
}

// Exec is a proxy for sql.DB.Exec.
func (d *DBInstance) Exec(query string, args ...interface{}) (ResultInstanceInterface, error) {
	res, err := d.FieldDB.Exec(query, args...)
	return &ResultInstance{FieldResult: res}, err
}

// Prepare is a proxy for sql.DB.Prepare.
func (d *DBInstance) Prepare(query string) (StmtInstanceInterface, error) {
	stmt, err := d.FieldDB.Prepare(query)
	return &StmtInstance{FieldStmt: stmt}, err
}

// Query is a proxy for sql.DB.Query.
func (d *DBInstance) Query(query string, args ...interface{}) (RowsInstanceInterface, error) {
	rows, _ := d.FieldDB.Query(query, args...)
	return &RowsInstance{rows}, nil
}
