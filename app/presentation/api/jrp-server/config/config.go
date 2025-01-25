package config

import (
	"path/filepath"
	"strings"

	baseConfig "github.com/yanosea/jrp/v2/app/config"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

// JrpServerConfigurator is an interface that gets the configuration of the Jrp server application.
type JrpServerConfigurator interface {
	GetConfig() (*JrpServerConfig, error)
}

// ServerConfigurator is a struct that implements the JrpServerConfigurator interface.
type ServerConfigurator struct {
	*baseConfig.BaseConfigurator
}

// NewJrpServerConfigurator creates a new JrpServerConfigurator.
func NewJrpServerConfigurator(
	envconfigProxy proxy.Envconfig,
	fileUtil utility.FileUtil,
) JrpServerConfigurator {
	return &ServerConfigurator{
		BaseConfigurator: baseConfig.NewConfigurator(
			envconfigProxy,
			fileUtil,
		),
	}
}

// JrpServerConfig is a struct that contains the configuration of the Jrp server application.
type JrpServerConfig struct {
	baseConfig.JrpConfig
	JrpPort string
}

// envConfig is a struct that contains the environment variables.
type envConfig struct {
	JrpPort     string          `envconfig:"JRP_SERVER_PORT" default:"8080"`
	WnJpnDBType database.DBType `envconfig:"JRP_SERVER_WNJPN_DB_TYPE" default:"sqlite"`
	WnJpnDBDsn  string          `envconfig:"JRP_SERVER_WNJPN_DB" default:"XDG_DATA_HOME/jrp/wnjpn.db"`
}

// GetConfig gets the configuration of the Jrp server application.
func (c *ServerConfigurator) GetConfig() (*JrpServerConfig, error) {
	var env envConfig
	if err := c.Envconfig.Process("", &env); err != nil {
		return nil, err
	}

	config := &JrpServerConfig{
		JrpConfig: baseConfig.JrpConfig{
			WNJpnDBType: env.WnJpnDBType,
			WNJpnDBDsn:  env.WnJpnDBDsn,
		},
		JrpPort: env.JrpPort,
	}

	if config.WNJpnDBType == database.SQLite {
		xdgDataHome, err := c.FileUtil.GetXDGDataHome()
		if err != nil {
			return nil, err
		}

		config.WNJpnDBDsn = strings.Replace(
			config.WNJpnDBDsn,
			"XDG_DATA_HOME",
			xdgDataHome,
			1,
		)
		if err := c.FileUtil.MkdirIfNotExist(
			filepath.Dir(config.WNJpnDBDsn),
		); err != nil {
			return nil, err
		}
	}

	return config, nil
}
