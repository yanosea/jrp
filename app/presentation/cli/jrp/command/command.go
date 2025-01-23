package command

import (
	"context"
	"os"

	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/config"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

var (
	// output is the output string
	output = ""
	// NewCli is a variable holding the current Cli creation function
	NewCli CreateCliFunc = newCli
)

type Cli interface {
	Init(envconfig proxy.Envconfig, sql proxy.Sql, version string, fileUtil utility.FileUtil, versionUtil utility.VersionUtil) int
	Run(ctx context.Context) int
}

// cli is a struct that represents the command line interface of jrp cli.
type cli struct {
	Cobra             proxy.Cobra
	Version           string
	RootCommand       proxy.Command
	ConnectionManager database.ConnectionManager
}

// CreateCliFunc is a function type for creating new Cli instances
type CreateCliFunc func(cobra proxy.Cobra) Cli

// newCli is the default implementation of CreateCliFunc
func newCli(cobra proxy.Cobra) Cli {
	return &cli{
		Cobra:       cobra,
		Version:     "",
		RootCommand: nil,
	}
}

// Init initializes the command line interface of jrp.
func (c *cli) Init(
	envconfig proxy.Envconfig,
	sql proxy.Sql,
	version string,
	fileUtil utility.FileUtil,
	versionUtil utility.VersionUtil,
) int {
	configurator := config.NewJrpCliConfigurator(envconfig, fileUtil)
	conf, err := configurator.GetConfig()
	if err != nil {
		output = formatter.AppendErrorToOutput(err, output)
		presenter.Print(os.Stderr, output)
		return 1
	}

	if c.ConnectionManager == nil {
		c.ConnectionManager = database.NewConnectionManager(sql)
	}

	if conf.JrpDBType == database.SQLite {
		if err := c.ConnectionManager.InitializeConnection(
			database.ConnectionConfig{
				DBName: database.JrpDB,
				DBType: conf.JrpDBType,
				DSN:    conf.JrpDBDsn,
			},
		); err != nil {
			output = formatter.AppendErrorToOutput(err, output)
			presenter.Print(os.Stderr, output)
			return 1
		}
	}

	if conf.WNJpnDBType == database.SQLite && fileUtil.IsExist(conf.WNJpnDBDsn) {
		if err := c.ConnectionManager.InitializeConnection(
			database.ConnectionConfig{
				DBName: database.WNJpnDB,
				DBType: conf.WNJpnDBType,
				DSN:    conf.WNJpnDBDsn,
			},
		); err != nil {
			output = formatter.AppendErrorToOutput(err, output)
			presenter.Print(os.Stderr, output)
			return 1
		}
	}

	ver := versionUtil.GetVersion(version)

	c.RootCommand = NewRootCommand(
		c.Cobra,
		ver,
		conf,
		&output,
	)

	return 0
}

// Run runs the command line interface of jrp cli.
func (c *cli) Run(ctx context.Context) (exitCode int) {
	defer func() {
		if c.ConnectionManager != nil {
			if err := c.ConnectionManager.CloseAllConnections(); err != nil {
				output = formatter.AppendErrorToOutput(err, output)
				presenter.Print(os.Stderr, output)
				exitCode = 1
			}
		}
	}()

	out := os.Stdout
	if err := c.RootCommand.ExecuteContext(ctx); err != nil {
		output = formatter.AppendErrorToOutput(err, output)
		out = os.Stderr
		exitCode = 1
	}

	presenter.Print(out, output)

	return
}
