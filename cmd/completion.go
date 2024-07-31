package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/util"
)

func newCompletionCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constant.COMPLETION_USE,
		Short: constant.COMPLETION_SHORT,
		Long:  constant.COMPLETION_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if no sub command is specified, print the message and return nil.
			util.PrintlnWithWriter(globalOption.Out, constant.COMPLETION_MESSAGE_NO_SUB_COMMAND)

			return nil
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_HELP_TEMPLATE)

	cmd.AddCommand(
		newCompletionBashCommand(globalOption),
		newCompletionFishCommand(globalOption),
		newCompletionPowerShellCommand(globalOption),
		newCompletionZshCommand(globalOption),
	)

	return cmd
}
