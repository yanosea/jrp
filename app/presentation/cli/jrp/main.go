package main

import (
	"context"
	"os"

	"github.com/yanosea/jrp/app/presentation/cli/jrp/command"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"
)

// JrpCliParams is a struct that represents the options of jrp cli.
type JrpCliParams struct {
	// Version is the version of jrp cli.
	Version string
	// Cobra is a proxy of spf13/cobra.
	Cobra proxy.Cobra
	// Envconfig is a proxy of kelseyhightower/envconfig.
	Envconfig proxy.Envconfig
	// Sql is a proxy of database/sql.
	Sql proxy.Sql
	// FileUtil provides the file utility.
	FileUtil utility.FileUtil
	// VersionUtil provides the version of the application.
	VersionUtil utility.VersionUtil
}

var (
	// version is the version of jrp cli and is embedded by goreleaser.
	version = ""
	// exit is a variable that contains the os.Exit function for injecting dependencies in testing.
	exit = os.Exit
	// jrpCliParams is a variable that contains the JrpCliParams struct.
	jrpCliParams = JrpCliParams{
		Version:   version,
		Cobra:     proxy.NewCobra(),
		Envconfig: proxy.NewEnvconfig(),
		Sql:       proxy.NewSql(),
		FileUtil: utility.NewFileUtil(
			proxy.NewGzip(),
			proxy.NewIo(),
			proxy.NewOs(),
		),
		VersionUtil: utility.NewVersionUtil(proxy.NewDebug()),
	}
)

// main is the entry point of jrp cli.
func main() {
	cli := command.NewCli(
		jrpCliParams.Cobra,
	)
	if exitCode := cli.Init(
		jrpCliParams.Envconfig,
		jrpCliParams.Sql,
		jrpCliParams.Version,
		jrpCliParams.FileUtil,
		jrpCliParams.VersionUtil,
	); exitCode != 0 {
		exit(exitCode)
	}

	ctx := context.Background()
	defer ctx.Done()

	exit(cli.Run(ctx))
}
