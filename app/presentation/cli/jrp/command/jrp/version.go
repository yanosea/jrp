package jrp

import (
	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

var (
	format = "plain"
)

// NewVersionCommand returns a new instance of the version command.
func NewVersionCommand(
	cobra proxy.Cobra,
	version string,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("version")
	cmd.SetAliases([]string{"ver", "v"})
	cmd.SetUsageTemplate(versionUsageTemplate)
	cmd.SetHelpTemplate(versionHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(_ *c.Command, _ []string) error {
			return runVersion(version, output)
		},
	)

	return cmd
}

// runVersion runs the version command.
func runVersion(version string, output *string) error {
	uc := jrpApp.NewGetVersionUseCase()
	dto := uc.Run(version)

	f, err := formatter.NewFormatter(format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o := f.Format(dto)
	*output = o

	return nil
}

const (
	// versionHelpTemplate is the help template of the version command.
	versionHelpTemplate = `🔖 Show the version of jrp

` + versionUsageTemplate
	// versionUsageTemplate is the usage template of the version command.
	versionUsageTemplate = `Usage:
  jrp version [flags]
	jrp ver     [flags]
	jrp v       [flags]

Flags:
  -h, --help  🤝 help for jrp version
`
)
