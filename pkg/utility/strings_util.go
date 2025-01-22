package utility

import (
	"strings"
)

// StringsUtil is an interface that contains the utility functions for manipulating strings.
type StringsUtil interface {
	RemoveNewLines(s string) string
	RemoveSpaces(s string) string
	RemoveTabs(s string) string
}

// stringsUtil is a struct that contains the utility functions for manipulating strings.
type stringsUtil struct{}

// NewStringsUtil returns a new instance of the StringsUtil struct.
func NewStringsUtil() StringsUtil {
	return &stringsUtil{}
}

// RemoveNewLines removes all new lines from the given strings.
func (s *stringsUtil) RemoveNewLines(str string) string {
	return strings.ReplaceAll(str, "\n", "")
}

// RemoveSpaces removes all spaces from the given strings.
func (s *stringsUtil) RemoveSpaces(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

// RemoveTabs removes all tabs from the given strings.
func (s *stringsUtil) RemoveTabs(str string) string {
	return strings.ReplaceAll(str, "\t", "")
}
