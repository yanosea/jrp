package cobraproxy

import (
	"github.com/spf13/cobra"
)

// Cobra is an interface for cobra.
type Cobra interface {
	MaximumNArgs(int) *PositionalArgsInstance
	NewCommand() *CommandInstance
}

// CobraProxy is a struct that implements Cobra.
type CobraProxy struct{}

// New is a constructor for CommandProxy.
func New() Cobra {
	return &CobraProxy{}
}

// MaximumNArgs is a proxy for cobra.MaximumNArgs.
func (*CobraProxy) MaximumNArgs(n int) *PositionalArgsInstance {
	return &PositionalArgsInstance{FieldPositionalArgs: cobra.MaximumNArgs(n)}
}

// NewCommand is a proxy for getting cobra.Command struct.
func (*CobraProxy) NewCommand() *CommandInstance {
	return &CommandInstance{FieldCommand: &cobra.Command{}}
}
