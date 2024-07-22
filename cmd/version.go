package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/util"
)

const (
	version_help_template = `🔖 Show the version of jrp.

Simply show the version of jrp.

Usage:
  jrp version [flags]

Flags:
  -h, --help   help for version
`
	version_use   = "version"
	version_short = "🔖 Show the version of jrp."
	version_long  = `🔖 Show the version of jrp.

Simply show the version of jrp.`
	version_message_template = "🔖 jrp version %s"
)

func newVersionCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   version_use,
		Short: version_short,
		Long:  version_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return globalOption.version()
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	cmd.SetHelpTemplate(version_help_template)

	return cmd
}

func (g *GlobalOption) version() error {
	// show version
	util.PrintlnWithWriter(g.Out, fmt.Sprintf(version_message_template, version))
	return nil
}
