package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
)

func newCompletionZshCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constant.COMPLETION_ZSH_USE,
		Short: constant.COMPLETION_ZSH_SHORT,
		Long:  constant.COMPLETION_ZSH_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenZshCompletion(globalOption.Out)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_ZSH_HELP_TEMPLATE)

	return cmd
}
