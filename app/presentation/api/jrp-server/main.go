package main

// @title JRP API
// @description jrp api server
// @host localhost:8080
// @BasePath /api

import (
	"os"

	"github.com/yanosea/jrp/v2/app/presentation/api/jrp-server/server"
	_ "github.com/yanosea/jrp/v2/docs"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

// JrpApiServerParams is a struct that represents the options of jrp api sever.
type JrpApiServerParams struct {
	// Echos is a proxy of labstack/echo/v4.
	Echos proxy.Echos
	// Envconfig is a proxy of kelseyhightower/envconfig.
	Envconfig proxy.Envconfig
	// FileUtil provides the file utility.
	FileUtil utility.FileUtil
	// Sql is a proxy of database/sql.
	Sql proxy.Sql
}

var (
	// exit is a variable that contains the os.Exit function for injecting dependencies in testing.
	exit = os.Exit
	// jrpApiServerParams is a variable that contains the jrpApiServerParams struct.
	jrpApiServerParams = JrpApiServerParams{
		Echos:     proxy.NewEchos(),
		Envconfig: proxy.NewEnvconfig(),
		FileUtil: utility.NewFileUtil(
			proxy.NewGzip(),
			proxy.NewIo(),
			proxy.NewOs(),
		),
		Sql: proxy.NewSql(),
	}
)

// main is the entry point of jrp api server.
func main() {
	serv := server.NewServer(
		jrpApiServerParams.Echos,
	)
	if exitCode := serv.Init(
		jrpApiServerParams.Envconfig,
		jrpApiServerParams.FileUtil,
		jrpApiServerParams.Sql,
	); exitCode != 0 {
		exit(exitCode)
	}

	exit(serv.Run())
}
