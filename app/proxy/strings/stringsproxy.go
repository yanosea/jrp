package stringsproxy

import (
	"strings"
)

// Strings is an interface for strings.
type Strings interface {
	Join(elems []string, sep string) string
}

// StringsProxy is a struct that implements Strings.
type StringsProxy struct{}

// New is a constructor for StringsProxy.
func New() Strings {
	return &StringsProxy{}
}

// Join is a proxy for strings.Join.
func (*StringsProxy) Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}
