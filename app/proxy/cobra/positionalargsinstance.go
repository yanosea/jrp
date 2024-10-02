package cobraproxy

import (
	"github.com/spf13/cobra"
)

// PositionalArgsInstanceInterface is an interface for cobra.PositionalArgs.
type PositionalArgsInstanceInterface interface {
}

// PositionalArgsInstance is a struct that implements PositionalArgsInstanceInterface.
type PositionalArgsInstance struct {
	FieldPositionalArgs cobra.PositionalArgs
}
