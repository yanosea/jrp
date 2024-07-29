package logic

import (
	"path/filepath"

	"github.com/yanosea/jrp/constant"
)

func GetDBFileDirPath(e Env, u User) (string, error) {
	// get current user
	currentUser, err := u.Current()
	if err != nil {
		return "", err
	}
	// set default path ($XDG_DATA_HOME/jrp)
	dbFileDirPath := filepath.Join(currentUser.HomeDir, ".local", "share", "jrp")
	// get JRP_ENV_WORDNETJP_DIR
	wordNetJpDirPath := e.Get(constant.JRP_ENV_WORDNETJP_DIR)
	if wordNetJpDirPath != "" {
		dbFileDirPath = wordNetJpDirPath
	}

	return dbFileDirPath, nil
}
