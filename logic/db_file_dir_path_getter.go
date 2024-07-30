package logic

import (
	"path/filepath"

	"github.com/yanosea/jrp/constant"
)

type FileDirPathGetter interface {
	GetFileDirPath() error
}

type DBFileDirPathGetter struct {
	Env  Env
	User User
}

func NewDBFileDirPathGetter(env Env, user User) *DBFileDirPathGetter {
	return &DBFileDirPathGetter{
		Env:  env,
		User: user,
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
	wordNetJpDirPath := g.Env.Get(constant.JRP_ENV_WORDNETJP_DIR)
	if wordNetJpDirPath != "" {
		dbFileDirPath = wordNetJpDirPath
	}

	return dbFileDirPath, nil
}
