package fmtproxy

import (
	"fmt"

	"github.com/yanosea/jrp/app/proxy/io"
)

// Fmt is an interface for fmt.
type Fmt interface {
	Fprintf(w ioproxy.WriterInstanceInterface, format string, a ...any)
	Sprintf(format string, a ...any) string
}

// FmtProxy is a struct that implements Fmt.
type FmtProxy struct{}

// New is a constructor for FmtProxy.
func New() Fmt {
	return &FmtProxy{}
}

// Fprintf is a proxy for fmt.Fprintf.
func (*FmtProxy) Fprintf(w ioproxy.WriterInstanceInterface, format string, a ...any) {
	fmt.Fprintf(w, format, a...)
}

// Sprintf is a proxy for fmt.Sprintf.
func (*FmtProxy) Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
