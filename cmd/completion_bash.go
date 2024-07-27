package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
)

func newCompletionBashCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constant.COMPLETION_BASH_USE,
		Short: constant.COMPLETION_BASH_SHORT,
		Long:  constant.COMPLETION_BASH_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenBashCompletion(globalOption.Out)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_BASH_HELP_TEMPLATE)

	return cmd
}
