package env

import (
	"os"
)

type EnvironmentProvider interface {
	Get(key string) string
}

type OsEnvironment struct{}

func (OsEnvironment) Get(key string) string {
	return os.Getenv(key)
}
