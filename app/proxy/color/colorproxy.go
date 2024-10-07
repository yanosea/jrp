package colorproxy

import (
	"github.com/fatih/color"
)

// Color is an interface for color.
type Color interface {
	GreenString(format string, a ...interface{}) string
	RedString(format string, a ...interface{}) string
	YellowString(format string, a ...interface{}) string
}

// ColorProxy is a struct that implements Color.
type ColorProxy struct{}

// New is a constructor for BufferProxy.
func New() Color {
	return &ColorProxy{}
}

// GreenString is a proxy for color.GreenString.
func (*ColorProxy) GreenString(format string, a ...interface{}) string {
	return color.GreenString(format, a...)
}

// RedString is a proxy for color.RedString.
func (*ColorProxy) RedString(format string, a ...interface{}) string {
	return color.RedString(format, a...)
}

// YellowString is a proxy for color.YellowString.
func (*ColorProxy) YellowString(format string, a ...interface{}) string {
	return color.YellowString(format, a...)
}
