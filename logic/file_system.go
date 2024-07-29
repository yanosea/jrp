package logic

import (
	"os"
)

type FileSystem interface {
	Create(name string) (*os.File, error)
	Remove(name string) error
}

type OSFileSystem struct{}

func (OSFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (OSFileSystem) Remove(name string) error {
	return os.Remove(name)
}
