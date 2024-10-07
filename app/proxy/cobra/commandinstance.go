package cobraproxy

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/pflag"
)

// CommandInstanceInterface is an interface for cobra.Command.
type CommandInstanceInterface interface {
	AddCommand(cmds ...*CommandInstance)
	Execute() error
	GetCommand() *cobra.Command
	PersistentFlags() *pflagproxy.FlagSetInstance
	SetArgs(args []string)
	SetErr(io ioproxy.WriterInstanceInterface)
	SetHelpTemplate(s string)
	SetOut(io ioproxy.WriterInstanceInterface)
}

// CommandInstance is a struct that implements CommandInstanceInterface.
type CommandInstance struct {
	FieldCommand *cobra.Command
}

// AddCommand is a proxy for cobra.Command.AddCommand.
func (c *CommandInstance) AddCommand(cmds ...*CommandInstance) {
	for _, cmd := range cmds {
		c.FieldCommand.AddCommand(cmd.FieldCommand)
	}
}

// Execute is a proxy for cobra.Command.Execute.
func (c *CommandInstance) Execute() error {
	return c.FieldCommand.Execute()
}

// GetCommand returns the cobra.Command.
func (c *CommandInstance) GetCommand() *cobra.Command {
	return c.FieldCommand
}

// PersistentFlags is a proxy for cobra.Command.PersistentFlags.
func (c *CommandInstance) PersistentFlags() *pflagproxy.FlagSetInstance {
	return &pflagproxy.FlagSetInstance{FieldFlagSet: c.FieldCommand.PersistentFlags()}
}

// SetArgs is a proxy for cobra.Command.SetArgs.
func (c *CommandInstance) SetArgs(args []string) {
	c.FieldCommand.SetArgs(args)
}

// SetErr is a proxy for cobra.Command.SetErr.
func (c *CommandInstance) SetErr(io ioproxy.WriterInstanceInterface) {
	c.FieldCommand.SetErr(io)
}

// SetHelpTemplate is a proxy for cobra.Command.SetHelpTemplate.
func (c *CommandInstance) SetHelpTemplate(s string) {
	c.FieldCommand.SetHelpTemplate(s)
}

// SetOut is a proxy for cobra.Command.SetOut.

func (c *CommandInstance) SetOut(io ioproxy.WriterInstanceInterface) {
	c.FieldCommand.SetOut(io)
}
