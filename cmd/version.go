package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/util"
)

func newVersionCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constant.VERSION_USE,
		Short: constant.VERSION_SHORT,
		Long:  constant.VERSION_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			return globalOption.version()
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)
	cmd.SetHelpTemplate(constant.VERSION_HELP_TEMPLATE)

	return cmd
}

func (g *GlobalOption) version() error {
	// show version
	util.PrintlnWithWriter(g.Out, fmt.Sprintf(constant.VERSION_MESSAGE_TEMPLATE, version))
	return nil
}
