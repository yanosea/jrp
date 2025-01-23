package jrp

import (
	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/config"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// NewDownloadCommand returns a new instance of the download command.
func NewDownloadCommand(
	cobra proxy.Cobra,
	conf *config.JrpCliConfig,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("download")
	cmd.SetAliases([]string{"dl", "d"})
	cmd.SetUsageTemplate(downloadUsageTemplate)
	cmd.SetHelpTemplate(downloadHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(_ *c.Command, _ []string) error {
			return runDownload(conf, output)
		},
	)

	return cmd
}

// runDownload runs the download command.
func runDownload(conf *config.JrpCliConfig, output *string) error {
	if conf.WNJpnDBType != database.SQLite {
		o := formatter.Red("❌ The type of WordNet Japan database is not sqlite...")
		*output = o
		return nil
	}

	if err := presenter.StartSpinner(
		true,
		"yellow",
		formatter.Yellow(
			"  📦 Downloading WordNet Japan sqlite database file from the official web site...",
		),
	); err != nil {
		o := formatter.Red("❌ Failed to start spinner...")
		*output = o
		return err
	}
	defer func() {
		presenter.StopSpinner()
	}()

	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(conf.WNJpnDBDsn); err != nil && err.Error() == "wnjpn.db already exists" {
		o := formatter.Green("✅ You are already ready to use jrp!")
		*output = o
		return nil
	} else if err != nil {
		o := formatter.Red("❌ Failed to download WordNet Japan sqlite database file...")
		*output = o
		return err
	}

	o := formatter.Green("✅ Downloaded successfully! Now, you are ready to use jrp!")
	*output = o

	return nil
}

const (
	// downloadHelpTemplate is the help template of the download command.
	downloadHelpTemplate = `📦 Download WordNet Japan sqlite database file from the official web site.

You have to download WordNet Japan sqlite database file to use jrp at first.
"jrp download" will download archive file from the official web site and decompress it to the database file.

The default directory is "$XDG_DATA_HOME/jrp" ("~/.local/share/jrp").
If you want to change the directory, set the “JRP_WNJPN_DB_FILE_DIR” environment variable.
You have to set the same directory to the “JRP_WNJPN_DB_FILE_DIR” environment variable when you use jrp.

` + downloadUsageTemplate
	// downloadUsageTemplate is the usage template of the download command.
	downloadUsageTemplate = `Usage:
  jrp download [flags]
  jrp dl       [flags]
  jrp d        [flags]

Flags:
  -h, --help  🤝 help for jrp download
`
)
