package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/logic"
)

type downloadOption struct {
	Out    io.Writer
	ErrOut io.Writer

	Env        logic.Env
	User       logic.User
	FileSystem logic.FileSystem
	HttpClient logic.HttpClient
}

func newDownloadCommand(globalOption *GlobalOption) *cobra.Command {
	o := &downloadOption{}
	cmd := &cobra.Command{
		Use:     constant.DOWNLOAD_USE,
		Aliases: constant.GetDownloadAliases(),
		Short:   constant.DOWNLOAD_SHORT,
		Long:    constant.DOWNLOAD_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut
			o.Env = logic.OsEnv{}
			o.User = logic.OsUser{}
			o.FileSystem = logic.OSFileSystem{}
			o.HttpClient = logic.DefaultHttpClient{}
			return o.download()
		},
	}

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.DOWNLOAD_HELP_TEMPLATE)

	return cmd
}

func (o *downloadOption) download() error {
	if err := logic.Download(o.Env, o.User, o.FileSystem, o.HttpClient); err != nil {
		return err
	}
	return nil
}
