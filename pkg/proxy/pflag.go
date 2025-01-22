package proxy

import (
	"github.com/spf13/pflag"
)

// FlagSet is an interface that provides a proxy of the methods of pflag.FlagSet.
type FlagSet interface {
	BoolVarP(p *bool, name string, shorthand string, value bool, usage string)
	IntVarP(p *int, name string, shorthand string, value int, usage string)
	StringVarP(p *string, name string, shorthand string, value string, usage string)
}

// flagSetProxy is a proxy struct that implements the FlagSet interface.
type flagSetProxy struct {
	flagSet *pflag.FlagSet
}

// BoolVarP returns a new instance of the FlagSet interface.
func (f *flagSetProxy) BoolVarP(p *bool, name string, shorthand string, value bool, usage string) {
	f.flagSet.BoolVarP(p, name, shorthand, value, usage)
}

// IntVarP returns a new instance of the FlagSet interface.
func (f *flagSetProxy) IntVarP(p *int, name string, shorthand string, value int, usage string) {
	f.flagSet.IntVarP(p, name, shorthand, value, usage)
}

// NewFlagSet returns a new instance of the FlagSet interface.
func (f *flagSetProxy) StringVarP(p *string, name string, shorthand string, value string, usage string) {
	f.flagSet.StringVarP(p, name, shorthand, value, usage)
}
