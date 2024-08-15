package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
)

func newCompletionPowerShellCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constant.COMPLETION_POWERSHELL_USE,
		Short: constant.COMPLETION_POWERSHELL_SHORT,
		Long:  constant.COMPLETION_POWERSHELL_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenPowerShellCompletion(globalOption.Out)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)
	cmd.SetHelpTemplate(constant.COMPLETION_POWERSHELL_HELP_TEMPLATE)

	return cmd
}
