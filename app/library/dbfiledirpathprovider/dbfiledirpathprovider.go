package dbfiledirpathprovider

import (
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/user"
)

// DBFileDirPathProvidable is an interface for DBFileDirPathProvider.
type DBFileDirPathProvidable interface {
	GetJrpDBFileDirPath() (string, error)
	GetWNJpnDBFileDirPath() (string, error)
}

// DBFileDirPathProvider is a struct that implements DBFileDirPathProvidable.
type DBFileDirPathProvider struct {
	FilepathProxy filepathproxy.FilePath
	OsProxy       osproxy.Os
	UserProxy     userproxy.User
}

// New is a constructor for DBFileDirPathProvider.
func New(
	filepath filepathproxy.FilePath,
	os osproxy.Os,
	user userproxy.User,
) *DBFileDirPathProvider {
	return &DBFileDirPathProvider{
		FilepathProxy: filepath,
		OsProxy:       os,
		UserProxy:     user,
	}
}

// GetJrpDBFileDirPath provides db file directory path for jrp db file.
func (d *DBFileDirPathProvider) GetJrpDBFileDirPath() (string, error) {
	return d.getDBFileDirPath(JRP_ENV_JRP_DB_FILE_DIR)
}

// GetWNJpnDBFileDirPath provides db file directory path for wnjpn db file.
func (d *DBFileDirPathProvider) GetWNJpnDBFileDirPath() (string, error) {
	return d.getDBFileDirPath(JRP_ENV_WNJPN_DB_FILE_DIR)
}

// getDBFileDirPath gets db file directory path from env var or default.
func (d *DBFileDirPathProvider) getDBFileDirPath(envVar string) (string, error) {
	// get env var
	envDir := d.OsProxy.Getenv(envVar)
	if envDir != "" {
		// if env var is set, use it
		return envDir, nil
	}

	// get current user
	currentUser, err := d.UserProxy.Current()
	if err != nil {
		return "", err
	}

	// get default db file dir path
	var defaultDBFileDirPath string
	xdgDataHome := d.OsProxy.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		// if XDG_DATA_HOME is set, use it
		defaultDBFileDirPath = d.FilepathProxy.Join(xdgDataHome, "jrp")
	} else {
		// if XDG_DATA_HOME is not set, use default
		defaultDBFileDirPath = d.FilepathProxy.Join(currentUser.FieldUser.HomeDir, ".local", "share", "jrp")
	}

	return defaultDBFileDirPath, nil
}
