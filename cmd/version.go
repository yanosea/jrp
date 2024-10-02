package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/library/versionprovider"
	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/debug"
	fmtproxy "github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/cmd/constant"
)

// NewVersionCommand creates a new version command.
func NewVersionCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.VERSION_USE
	cmd.FieldCommand.Short = constant.VERSION_SHORT
	cmd.FieldCommand.Long = constant.VERSION_LONG
	cmd.FieldCommand.RunE = g.versionRunE

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.VERSION_HELP_TEMPLATE)

	return cmd
}

// versionRunE is the function that is called when the version command is executed.
func (g *GlobalOption) versionRunE(_ *cobra.Command, _ []string) error {
	return g.version()
}

// version shows the version of jrp.
func (g *GlobalOption) version() error {
	v := versionprovider.New(debugproxy.New())
	fmtProxy := fmtproxy.New()
	// get version from buildinfo and write it
	g.Utility.PrintlnWithWriter(g.Out, fmtProxy.Sprintf(constant.VERSION_MESSAGE_TEMPLATE, v.GetVersion(ver)))

	return nil
}
