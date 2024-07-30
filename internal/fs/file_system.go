package fs

import (
	"os"
)

type FileManager interface {
	Create(name string) (*os.File, error)
	Remove(name string) error
}

type OsFileManager struct{}

func (OsFileManager) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (OsFileManager) Remove(name string) error {
	return os.Remove(name)
}
