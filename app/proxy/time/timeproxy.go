package timeproxy

import (
	"time"
)

// Time is an interface for time.
type Time interface {
	Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) *TimeInstance
	Now() *TimeInstance
}

// TimeProxy is a struct that implements Time.
type TimeProxy struct{}

// New is a constructor for TimeProxy.
func New() Time {
	return &TimeProxy{}
}

// Date is a proxy for time.Date.
func (t *TimeProxy) Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) *TimeInstance {
	return &TimeInstance{FieldTime: time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

// Now is a proxy for time.Now.
func (*TimeProxy) Now() *TimeInstance {
	return &TimeInstance{FieldTime: time.Now()}
}
