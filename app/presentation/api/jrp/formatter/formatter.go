package formatter

import (
	"errors"
)

// Formatter is an interface that formats the output of jrp cli.
type Formatter interface {
	Format(result interface{}) ([]byte, error)
}

// NewFormatter returns a new instance of the Formatter interface.
func NewFormatter(
	format string,
) (Formatter, error) {
	var f Formatter
	switch format {
	case "json":
		f = NewJsonFormatter()
	default:
		return nil, errors.New("invalid format")
	}
	return f, nil
}
