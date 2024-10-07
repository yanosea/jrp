package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/cmd/constant"
)

// NewCompletionBashCommand creates a new completion bash command.
func NewCompletionBashCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.COMPLETION_BASH_USE
	cmd.FieldCommand.RunE = g.completionBashRunE

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_BASH_HELP_TEMPLATE)

	return cmd
}

// completionBashRunE is a function that is called when the completion bash command is executed.
func (g *GlobalOption) completionBashRunE(c *cobra.Command, _ []string) error {
	return g.completionBash(c)
}

// completionBash generates the bash completion script.
func (g *GlobalOption) completionBash(c *cobra.Command) error {
	return c.GenBashCompletion(g.Out)
}
