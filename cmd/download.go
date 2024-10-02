package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/downloader"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/gzip"
	"github.com/yanosea/jrp/app/proxy/http"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/spinner"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// downloadOption is the struct for download command.
type downloadOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	Downloader            downloader.Downloadable
	SpinnerProxy          spinnerproxy.Spinner
	Utility               utility.UtilityInterface
}

// NewDownloadCommand creates a new download command.
func NewDownloadCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &downloadOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Utility: g.Utility,
	}
	o.DBFileDirPathProvider = dbfiledirpathprovider.New(
		filepathproxy.New(),
		osproxy.New(),
		userproxy.New(),
	)
	o.Downloader = downloader.New(
		filepathproxy.New(),
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osproxy.New(),
		g.Utility,
	)
	o.SpinnerProxy = spinnerproxy.New()

	cobraproxy := cobraproxy.New()
	cmd := cobraproxy.NewCommand()

	cmd.FieldCommand.Use = constant.DOWNLOAD_USE
	cmd.FieldCommand.Aliases = constant.GetDownloadAliases()
	cmd.FieldCommand.Short = constant.DOWNLOAD_SHORT
	cmd.FieldCommand.Long = constant.DOWNLOAD_LONG
	cmd.FieldCommand.RunE = o.downloadRunE

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.DOWNLOAD_HELP_TEMPLATE)

	return cmd
}

// downloadRunE is the function that is called when the download command is executed.
func (o *downloadOption) downloadRunE(_ *cobra.Command, _ []string) error {
	// get wnjpn db file dir path
	wnJpnDBFileDirPath, err := o.DBFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		return err
	}

	// create the directory if it does not exist
	if err := o.Utility.CreateDirIfNotExist(wnJpnDBFileDirPath); err != nil {
		return err
	}

	return o.download(wnJpnDBFileDirPath)
}

// download downloads wnjpn db file.
func (o *downloadOption) download(wnJpnDBFileDirPath string) error {
	// start spinner
	spinner := o.SpinnerProxy.NewSpinner()
	spinner.Reverse()
	if err := spinner.SetColor("yellow"); err != nil {
		return err
	}
	colorProxy := colorproxy.New()
	spinner.SetSuffix(colorProxy.YellowString(constant.DOWNLOAD_MESSAGE_DOWNLOADING))
	spinner.Start()

	// download with downloader
	res, err := o.Downloader.DownloadWNJpnDBFile(wnJpnDBFileDirPath)
	spinner.Stop()
	o.writeDownloadResult(res)

	return err
}

// writeDownloadResult writes the download result.
func (o *downloadOption) writeDownloadResult(result downloader.DownloadStatus) {
	colorProxy := colorproxy.New()
	if result == downloader.DownloadedFailed {
		o.Utility.PrintlnWithWriter(o.ErrOut, colorProxy.RedString(constant.DOWNLOAD_MESSAGE_FAILED))
	} else if result == downloader.DownloadedAlready {
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.GreenString(constant.DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED))
	} else {
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.GreenString(constant.DOWNLOAD_MESSAGE_SUCCEEDED))
	}
}
