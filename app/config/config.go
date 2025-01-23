package config

import (
	"github.com/yanosea/jrp/v2/app/infrastructure/database"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

// Configurator is an interface that gets the configuration.
type Configurator interface {
	GetConfig() (*JrpConfig, error)
}

// BaseConfigurator is a struct that implements the Configurator interface.
type BaseConfigurator struct {
	Envconfig proxy.Envconfig
	FileUtil  utility.FileUtil
}

// JrpConfig is a struct that contains the configuration of the Jrp application.
type JrpConfig struct {
	JrpDBType   database.DBType
	JrpDBDsn    string
	WNJpnDBType database.DBType
	WNJpnDBDsn  string
}

// GetConfig gets the configuration of the Jrp application.
func NewConfigurator(
	envconfigProxy proxy.Envconfig,
	fileUtil utility.FileUtil,
) *BaseConfigurator {
	return &BaseConfigurator{
		Envconfig: envconfigProxy,
		FileUtil:  fileUtil,
	}
}
