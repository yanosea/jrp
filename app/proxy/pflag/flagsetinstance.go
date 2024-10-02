package pflagproxy

import (
	"github.com/spf13/pflag"
)

// FlagSetInstanceInterface is an interface for pflag.FlagSet.
type FlagSetInstanceInterface interface {
	BoolVarP(p *bool, name string, shorthand string, value bool, usage string)
	IntVarP(p *int, name string, shorthand string, value int, usage string)
	StringVarP(p *string, name string, shorthand string, value string, usage string)
}

// FlagSetInstance is a struct that implements FlagSetInstanceInterface.
type FlagSetInstance struct {
	FieldFlagSet *pflag.FlagSet
}

// BoolVarP is a proxy for pflag.FlagSet.BoolVarP.
func (f *FlagSetInstance) BoolVarP(p *bool, name string, shorthand string, value bool, usage string) {
	f.FieldFlagSet.BoolVarP(p, name, shorthand, value, usage)
}

// IntVarP is a proxy for pflag.FlagSet.IntVarP.
func (f *FlagSetInstance) IntVarP(p *int, name string, shorthand string, value int, usage string) {
	f.FieldFlagSet.IntVarP(p, name, shorthand, value, usage)
}

// StringVarP is a proxy for pflag.FlagSet.StringVarP.
func (f *FlagSetInstance) StringVarP(p *string, name string, shorthand string, value string, usage string) {
	f.FieldFlagSet.StringVarP(p, name, shorthand, value, usage)
}
