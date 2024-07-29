package logic

import (
	"path/filepath"

	"github.com/yanosea/jrp/constant"
)

func GetDBFileDirPath(e Env, u User) (string, error) {
	// check if JRP_ENV_WORDNETJP_DIR is set
	dbFileDirPath := e.Get(constant.JRP_ENV_WORDNETJP_DIR)
	if dbFileDirPath == "" {
		// get current user
		currentUser, err := u.Current()
		if err != nil {
			return "", err
		}
		// set default path ($XDG_DATA_HOME/jrp)
		dbFileDirPath = filepath.Join(currentUser.HomeDir, ".local", "share", "jrp")
	}
	return dbFileDirPath, nil
}
