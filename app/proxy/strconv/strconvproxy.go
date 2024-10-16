package strconvproxy

import (
	"strconv"
)

// Strconv is an interface for strconv.
type Strconv interface {
	Atoi(s string) (int, error)
	FormatBool(b bool) string
	Itoa(i int) string
}

// StrconvProxy is a struct that implements Strconv.
type StrconvProxy struct{}

// New is a constructor for StrconvProxy.
func New() Strconv {
	return &StrconvProxy{}
}

// Atoi is a proxy for strconv.Atoi.
func (*StrconvProxy) Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

func (*StrconvProxy) FormatBool(b bool) string {
	return strconv.FormatBool(b)
}

// Itoa is a proxy for strconv.Itoa.
func (*StrconvProxy) Itoa(i int) string {
	return strconv.Itoa(i)
}
