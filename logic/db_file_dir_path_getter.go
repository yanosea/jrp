package logic

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/yanosea/jrp/constant"
)

type Env interface {
	Get(key string) string
}

type UserProvider interface {
	Current() (*user.User, error)
}

type OsEnv struct{}

func (o OsEnv) Get(key string) string {
	return os.Getenv(key)
}

type OsUser struct{}

func (o OsUser) Current() (*user.User, error) {
	return user.Current()
}

func GetDBFileDirPath(env Env, userProvider UserProvider) (string, error) {
	// check if JRP_ENV_WORDNETJP_DIR is set
	dbFileDirPath := env.Get(constant.JRP_ENV_WORDNETJP_DIR)
	if dbFileDirPath == "" {
		// get current user
		user, err := userProvider.Current()
		if err != nil {
			return "", err
		}
		// default path ($XDG_DATA_HOME/jrp)
		dbFileDirPath = filepath.Join(user.HomeDir, ".local", "share", "jrp")
	}
	return dbFileDirPath, nil
}
