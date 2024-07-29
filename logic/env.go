package logic

import (
	"os"
)

type Env interface {
	Get(key string) string
}

type OsEnv struct{}

func (e OsEnv) Get(key string) string {
	return os.Getenv(key)
}
