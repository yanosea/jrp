package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"
)

type downloadOption struct {
	Out    io.Writer
	ErrOut io.Writer
}

func newDownloadCommand(g *GlobalOption) *cobra.Command {
	o := &downloadOption{
		Out:    g.Out,
		ErrOut: g.ErrOut,
	}

	cmd := &cobra.Command{
		Use:     constant.DOWNLOAD_USE,
		Aliases: constant.GetDownloadAliases(),
		Short:   constant.DOWNLOAD_SHORT,
		Long:    constant.DOWNLOAD_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.download()
		},
	}

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.DOWNLOAD_HELP_TEMPLATE)

	return cmd
}

func (o *downloadOption) download() error {
	u := usermanager.OSUserProvider{}
	f := fs.OsFileManager{}
	h := httpclient.DefaultHTTPClient{}
	i := iomanager.DefaultIOHelper{}
	g := gzip.DefaultGzipHandler{}

	downloader := logic.NewDBFileDownloader(u, f, h, i, g)

	if err := downloader.Download(); err != nil {
		return err
	}
	return nil
}
