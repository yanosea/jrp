package formatter

import (
	"fmt"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
)

// PlainFormatter is a struct that formats the output of jrp cli.
type PlainFormatter struct{}

// NewPlainFormatter returns a new instance of the PlainFormatter struct.
func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

// Format formats the output of jrp cli.
func (f *PlainFormatter) Format(result interface{}) (string, error) {
	var formatted string
	switch v := result.(type) {
	case *jrpApp.GetVersionUseCaseOutputDto:
		formatted = fmt.Sprintf("jrp version %s", v.Version)
	case []*jrpApp.GenerateJrpUseCaseOutputDto:
		for i, item := range v {
			formatted += item.Phrase
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	case []*jrpApp.GetHistoryUseCaseOutputDto:
		for i, item := range v {
			formatted += item.Phrase
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	case []*jrpApp.SearchHistoryUseCaseOutputDto:
		for i, item := range v {
			formatted += item.Phrase
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	default:
		formatted = ""
	}
	return formatted, nil
}
