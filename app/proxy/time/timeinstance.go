package timeproxy

import (
	"database/sql/driver"
	"time"
)

// TimeInstanceInterface is an interface for time.Time.
type TimeInstanceInterface interface {
	Format(layout string) string
	Scan(value interface{}) error
	Value() (driver.Value, error)
}

// TimeInstance is a struct that implements TimeInstanceInterface.
type TimeInstance struct {
	FieldTime time.Time
}

// Format is a proxy for time.Time.Format.
func (t *TimeInstance) Format(layout string) string {
	return t.FieldTime.Format(layout)
}

// Scan implements the sql.Scanner interface.
func (t *TimeInstance) Scan(value interface{}) error {
	t.FieldTime = value.(time.Time)
	return nil
}

// Value implements the driver.Valuer interface.
func (t *TimeInstance) Value() (driver.Value, error) {
	return t.FieldTime, nil
}
