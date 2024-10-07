package timeproxy

import (
	"time"
)

const (
	// Millisecond is  a proxy for time.Millisecond.
	Millisecond = time.Millisecond
)

// UTC is a variable for time.UTC.
var UTC = *time.UTC
