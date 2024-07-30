package fs

import (
	"os"
)

type FileManager interface {
	Create(name string) (*os.File, error)
	Remove(name string) error
	Exists(path string) bool
}

type OsFileManager struct{}

func (OsFileManager) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (OsFileManager) Remove(name string) error {
	return os.Remove(name)
}

func (o OsFileManager) Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
