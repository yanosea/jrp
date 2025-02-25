package config

import (
	"path/filepath"
	"strings"

	baseConfig "github.com/yanosea/jrp/v2/app/config"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

// JrpCliConfigurator is an interface that gets the configuration of the Jrp cli application.
type JrpCliConfigurator interface {
	GetConfig() (*JrpCliConfig, error)
}

// cliConfigurator is a struct that implements the JrpCliConfigurator interface.
type cliConfigurator struct {
	*baseConfig.BaseConfigurator
}

// NewJrpCliConfigurator creates a new JrpCliConfigurator.
func NewJrpCliConfigurator(
	envconfigProxy proxy.Envconfig,
	fileUtil utility.FileUtil,
) JrpCliConfigurator {
	return &cliConfigurator{
		BaseConfigurator: baseConfig.NewConfigurator(
			envconfigProxy,
			fileUtil,
		),
	}
}

// JrpCliConfig is a struct that contains the configuration of the Jrp cli application.
type JrpCliConfig struct {
	baseConfig.JrpConfig
	JrpDBType database.DBType
	JrpDBDsn  string
}

// envConfig is a struct that contains the environment variables.
type envConfig struct {
	JrpDBType   database.DBType `envconfig:"JRP_DB_TYPE" default:"sqlite"`
	JrpDBDsn    string          `envconfig:"JRP_DB" default:"XDG_DATA_HOME/jrp/jrp.db"`
	WnJpnDBType database.DBType `envconfig:"JRP_WNJPN_DB_TYPE" default:"sqlite"`
	WnJpnDBDsn  string          `envconfig:"JRP_WNJPN_DB" default:"XDG_DATA_HOME/jrp/wnjpn.db"`
}

// GetConfig gets the configuration of the Jrp cli application.
func (c *cliConfigurator) GetConfig() (*JrpCliConfig, error) {
	var env envConfig
	if err := c.Envconfig.Process("", &env); err != nil {
		return nil, err
	}

	config := &JrpCliConfig{
		JrpConfig: baseConfig.JrpConfig{
			WNJpnDBType: env.WnJpnDBType,
			WNJpnDBDsn:  env.WnJpnDBDsn,
		},
		JrpDBType: env.JrpDBType,
		JrpDBDsn:  env.JrpDBDsn,
	}

	if config.JrpDBType == database.SQLite || config.WNJpnDBType == database.SQLite {
		xdgDataHome, err := c.FileUtil.GetXDGDataHome()
		if err != nil {
			return nil, err
		}

		if config.JrpDBType == database.SQLite {
			config.JrpDBDsn = strings.Replace(
				config.JrpDBDsn,
				"XDG_DATA_HOME",
				xdgDataHome,
				1,
			)
			if err := c.FileUtil.MkdirIfNotExist(
				filepath.Dir(config.JrpDBDsn),
			); err != nil {
				return nil, err
			}
		}

		if config.WNJpnDBType == database.SQLite {
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
	}

	return config, nil
}
