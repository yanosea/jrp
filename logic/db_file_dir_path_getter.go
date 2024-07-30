package logic

import (
	"os"
	"path/filepath"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/usermanager"
)

type FileDirPathGetter interface {
	GetFileDirPath() error
}

type DBFileDirPathGetter struct {
	User usermanager.UserProvider
}

func NewDBFileDirPathGetter(u usermanager.UserProvider) *DBFileDirPathGetter {
	return &DBFileDirPathGetter{
		User: u,
	}
}

func (g *DBFileDirPathGetter) GetFileDirPath() (string, error) {
	// get current user
	currentUser, err := g.User.Current()
	if err != nil {
		return "", err
	}
	// set default path ($XDG_DATA_HOME/jrp)
	dbFileDirPath := filepath.Join(currentUser.HomeDir, ".local", "share", "jrp")
	// get JRP_ENV_WORDNETJP_DIR
	wordNetJpDirPath := os.Getenv(constant.JRP_ENV_WORDNETJP_DIR)
	if wordNetJpDirPath != "" {
		dbFileDirPath = wordNetJpDirPath
	}

	return dbFileDirPath, nil
}
