package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/cmd/constant"
)

// NewCompletionZshCommand creates a new completion zsh command.
func NewCompletionZshCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.COMPLETION_ZSH_USE
	cmd.FieldCommand.Short = constant.COMPLETION_ZSH_SHORT
	cmd.FieldCommand.Long = constant.COMPLETION_ZSH_LONG
	cmd.FieldCommand.RunE = g.completionZshRunE

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_ZSH_HELP_TEMPLATE)

	return cmd
}

// completionZshRunE is a function that is called when the completion zsh command is executed.
func (g *GlobalOption) completionZshRunE(c *cobra.Command, _ []string) error {
	return g.completionZsh(c)
}

// completionZsh generates the zsh completion script.
func (g *GlobalOption) completionZsh(c *cobra.Command) error {
	return c.GenZshCompletion(g.Out)
}
