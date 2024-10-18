package testutility

import (
	"strings"
)

const (
	// TEST_OUTPUT_STDOUT is a constant for stdout
	TEST_OUTPUT_ANY = "ANY"
)

// ReplaceDoubleSlashToSingleSlash replaces double slash to single slash
func ReplaceDoubleSlashToSingleSlash(input string) string {
	return strings.ReplaceAll(input, "//", "/")
}

// RemoveTabAndSpaceAndLf removes tab, space, and line feed
func RemoveTabAndSpaceAndLf(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "\t", ""), " ", ""), "\n", "")
}
