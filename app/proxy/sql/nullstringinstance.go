package sqlproxy

import (
	"database/sql"
	"database/sql/driver"
)

// NullStringInstanceInterface is an interface for sql.NullString.
type NullStringInstanceInterface interface {
	Scan(value interface{}) error
	Value() (driver.Value, error)
}

// NullStringInstance is a struct that implements NullStringInstanceInterface.
type NullStringInstance struct {
	FieldNullString *sql.NullString
}

// Scan implements the sql.Scanner interface.
func (n *NullStringInstance) Scan(value interface{}) error {
	if n.FieldNullString == nil {
		n.FieldNullString = &sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	str, _ := value.(string)
	n.FieldNullString.String, n.FieldNullString.Valid = str, true
	return nil
}

// Value implements the driver Valuer interface.
func (n *NullStringInstance) Value() (driver.Value, error) {
	var v driver.Value
	if n.FieldNullString == nil || !n.FieldNullString.Valid {
		v = nil
	}
	return v, nil
}
