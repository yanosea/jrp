package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/cmd/constant"
)

// NewCompletionCommand creates a new completion command.
func NewCompletionCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.COMPLETION_USE
	cmd.FieldCommand.RunE = g.completionRunE

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_HELP_TEMPLATE)

	cmd.AddCommand(
		NewCompletionBashCommand(g),
		NewCompletionFishCommand(g),
		NewCompletionPowerShellCommand(g),
		NewCompletionZshCommand(g),
	)

	return cmd
}

// completionRunE is the function that is called when the completion command is executed.
func (g *GlobalOption) completionRunE(_ *cobra.Command, _ []string) error {
	return g.completion()
}

// completion just prints the message.
func (g *GlobalOption) completion() error {
	// if no sub command is specified, print the message and return nil.
	colorProxy := colorproxy.New()
	g.Utility.PrintlnWithWriter(g.Out, colorProxy.YellowString(constant.COMPLETION_MESSAGE_NO_SUB_COMMAND))

	return nil
}
