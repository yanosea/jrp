package proxy

import (
	"context"

	"github.com/spf13/cobra"
)

// Cobra is an interface that provides a proxy of the methods of cobra
type Cobra interface {
	ExactArgs(n int) cobra.PositionalArgs
	MaximumNArgs(n int) cobra.PositionalArgs
	NewCommand() Command
}

// cobraProxy is a proxy struct that implements the Cobra interface.
type cobraProxy struct{}

// NewCobra returns a new instance of the Cobra interface.
func NewCobra() Cobra {
	return &cobraProxy{}
}

// ExactArgs is a proxy method that returns the cobra.ExactArgs.
func (*cobraProxy) ExactArgs(n int) cobra.PositionalArgs {
	return cobra.ExactArgs(n)
}

// MaximumArgs is a proxy method that returns the cobra.MaximumArgs.
func (*cobraProxy) MaximumNArgs(n int) cobra.PositionalArgs {
	return cobra.MaximumNArgs(n)
}

// NewCommand returns a new instance of the Command interface.
func (*cobraProxy) NewCommand() Command {
	return &commandProxy{command: &cobra.Command{}}
}

// Command is an interface that provides a proxy of the methods of cobra.Command.
type Command interface {
	AddCommand(cmds ...Command)
	ExecuteContext(ctx context.Context) error
	Flags() FlagSet
	GetCommand() *cobra.Command
	PersistentFlags() FlagSet
	RunE(cmd *cobra.Command, args []string) error
	SetAliases(s []string)
	SetArgs(f cobra.PositionalArgs)
	SetHelpTemplate(s string)
	SetRunE(f func(*cobra.Command, []string) error)
	SetSilenceErrors(b bool)
	SetUse(s string)
	SetUsageTemplate(s string)
}

// commandProxy is a proxy struct that implements the Command interface.
type commandProxy struct {
	command *cobra.Command
}

// AddCommand is a proxy method that calls the AddCommand method of the cobra.Command.
func (c *commandProxy) AddCommand(cmds ...Command) {
	for _, cmd := range cmds {
		c.command.AddCommand(cmd.GetCommand())
	}
}

// ExecuteContext is a proxy method that calls the ExecuteContext method of the cobra.Command.
func (c *commandProxy) ExecuteContext(ctx context.Context) error {
	return c.command.ExecuteContext(ctx)
}

// Flags is a proxy method that returns the cobra.FlagSet.
func (c *commandProxy) Flags() FlagSet {
	return &flagSetProxy{flagSet: c.command.Flags()}
}

// GetCommand is a proxy method that returns the cobra.Command.
func (c *commandProxy) GetCommand() *cobra.Command {
	return c.command
}

// PersistentFlags is a proxy method that returns the cobra.FlagSet.
func (c *commandProxy) PersistentFlags() FlagSet {
	return &flagSetProxy{flagSet: c.command.PersistentFlags()}
}

// RunE is a proxy method that calls the RunE method of the cobra.Command.
func (c *commandProxy) RunE(cmd *cobra.Command, args []string) error {
	return c.command.RunE(cmd, args)
}

// SetAlias is a proxy method that sets the Alias field of the cobra.Command.
func (c *commandProxy) SetAliases(s []string) {
	c.command.Aliases = s
}

// SetArgs is a proxy method that calls the SetArgs method of the cobra.Command.
func (c *commandProxy) SetArgs(p cobra.PositionalArgs) {
	c.command.Args = p
}

// SetHelpTemplate is a proxy method that calls the SetHelpTemplate method of the cobra.Command.
func (c *commandProxy) SetHelpTemplate(s string) {
	c.command.SetHelpTemplate(s)
}

// SetRunE is a proxy method that calls the SetRunE method of the cobra.Command.
func (c *commandProxy) SetRunE(f func(*cobra.Command, []string) error) {
	c.command.RunE = f
}

// SetSilenceErrors is a proxy method that sets the SilenceErrors field of the cobra.Command.
func (c *commandProxy) SetSilenceErrors(b bool) {
	c.command.SilenceErrors = b
}

// SetUse is a proxy method that sets the Use field of the cobra.Command.
func (c *commandProxy) SetUse(s string) {
	c.command.Use = s
}

// SetUsageTemplate is a proxy method that calls the SetUsageTemplate method of the cobra.Command.
func (c *commandProxy) SetUsageTemplate(s string) {
	c.command.SetUsageTemplate(s)
}
