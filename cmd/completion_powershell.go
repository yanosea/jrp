package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/cmd/constant"
)

// NewCompletionPowerShellCommand creates a new completion powershell command.
func NewCompletionPowerShellCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.COMPLETION_POWERSHELL_USE
	cmd.FieldCommand.RunE = g.completionPowerShellRunE

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_POWERSHELL_HELP_TEMPLATE)

	return cmd
}

// completionPowerShellRunE is a function that is called when the completion powershell command is executed.
func (g *GlobalOption) completionPowerShellRunE(c *cobra.Command, _ []string) error {
	return g.completionPowerShell(c)
}

// completionPowerShell generates the powershell completion script.
func (g *GlobalOption) completionPowerShell(c *cobra.Command) error {
	return c.GenPowerShellCompletion(g.Out)
}
