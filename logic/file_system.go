package logic

import (
	"os"
)

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (*os.File, error)
	Remove(name string) error
	TempDir() string
}

type OSFileSystem struct{}

func (OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OSFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (OSFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (OSFileSystem) TempDir() string {
	return os.TempDir()
}
