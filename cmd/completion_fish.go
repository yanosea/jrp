package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
)

func newCompletionFishCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constant.COMPLETION_FISH_USE,
		Short: constant.COMPLETION_FISH_SHORT,
		Long:  constant.COMPLETION_FISH_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenFishCompletion(globalOption.Out, false)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_FISH_HELP_TEMPLATE)

	return cmd
}
