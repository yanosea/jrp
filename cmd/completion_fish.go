package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/cmd/constant"
)

// NewCompletionFishCommand creates a new command for fish completion.
func NewCompletionFishCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.COMPLETION_FISH_USE
	cmd.FieldCommand.Short = constant.COMPLETION_FISH_SHORT
	cmd.FieldCommand.Long = constant.COMPLETION_FISH_LONG
	cmd.FieldCommand.RunE = g.completionFishRunE

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_FISH_HELP_TEMPLATE)

	return cmd
}

// completionFishRunE is the function that is called when the completion fish command is executed.
func (g *GlobalOption) completionFishRunE(c *cobra.Command, _ []string) error {
	return g.completionFish(c)
}

// completionFish generates the fish completion script.
func (g *GlobalOption) completionFish(c *cobra.Command) error {
	return c.GenFishCompletion(g.Out, false)
}
