package server

import (
	"errors"

	"github.com/labstack/echo/v4/middleware"

	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/presentation/api/jrp-server/config"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

var (
	// NewServer is a variable holding the current server creation function.
	NewServer CreateServerFunc = newServer
)

// Server is an interface that provides a proxy of the methods of jrp server.
type Server interface {
	Init(envconfig proxy.Envconfig, fileUtil utility.FileUtil, sql proxy.Sql) int
	Run() int
}

// server is a struct that represents the server interface of jrp server.
type server struct {
	ConnectionManager database.ConnectionManager
	Echos             proxy.Echos
	Logger            proxy.Logger
	Port              string
	Route             proxy.Echo
}

// CreateServerFunc is a function type for creating new server instances.
type CreateServerFunc func(echo proxy.Echos) Server

// newServer is the default implementation of CreateServerFunc.
func newServer(echos proxy.Echos) Server {
	return &server{
		ConnectionManager: nil,
		Echos:             echos,
		Logger:            nil,
		Port:              "",
		Route:             nil,
	}
}

// Init initializes the server.
func (s *server) Init(
	envconfig proxy.Envconfig,
	fileUtil utility.FileUtil,
	sql proxy.Sql,
) int {
	s.Route, s.Logger = s.Echos.NewEcho()
	s.Route.Use(middleware.Logger())
	s.Route.Use(middleware.Recover())
	Bind(s.Route)

	configurator := config.NewJrpServerConfigurator(envconfig, fileUtil)
	conf, err := configurator.GetConfig()
	if err != nil {
		s.Logger.Fatal(err)
		return 1
	}

	s.Port = conf.JrpPort

	if s.ConnectionManager == nil {
		s.ConnectionManager = database.NewConnectionManager(sql)
	}

	dbConfig := database.ConnectionConfig{
		DBName: database.WNJpnDB,
		DBType: conf.WNJpnDBType,
		DSN:    conf.WNJpnDBDsn,
	}

	if conf.WNJpnDBType == database.SQLite && !fileUtil.IsExist(conf.WNJpnDBDsn) {
		s.Logger.Fatal(errors.New("WordNet Japanese database is not found..."))
		return 1
	}

	if err := s.ConnectionManager.InitializeConnection(dbConfig); err != nil {
		s.Logger.Fatal(err)
		return 1
	}

	return 0
}

// Run runs the server.
func (s *server) Run() (exitCode int) {
	defer func() {
		if s.ConnectionManager != nil {
			if err := s.ConnectionManager.CloseAllConnections(); err != nil {
				s.Logger.Fatal(err)
				exitCode = 1
			}
		}
	}()

	if err := s.Route.Start(":" + s.Port); err != nil {
		s.Logger.Fatal(err)
		exitCode = 1
	}

	return
}
