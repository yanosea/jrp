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

type DownloadOption struct {
	Out        io.Writer
	ErrOut     io.Writer
	Downloader logic.Downloader
}

func NewDownloadCommand(g *GlobalOption) *cobra.Command {
	o := &DownloadOption{
		Out:    g.Out,
		ErrOut: g.ErrOut,
	}

	cmd := &cobra.Command{
		Use:     constant.DOWNLOAD_USE,
		Aliases: constant.GetDownloadAliases(),
		Short:   constant.DOWNLOAD_SHORT,
		Long:    constant.DOWNLOAD_LONG,
		RunE:    o.DownloadRunE,
	}

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.DOWNLOAD_HELP_TEMPLATE)

	return cmd
}

func (o *DownloadOption) DownloadRunE(_ *cobra.Command, _ []string) error {
	o.Downloader = logic.NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
	return o.Download()
}

func (o *DownloadOption) Download() error {
	if err := o.Downloader.Download(); err != nil {
		return err
	}
	return nil
}
