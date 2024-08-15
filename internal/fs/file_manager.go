package fs

import (
	"os"
)

type FileManager interface {
	Create(name string) (File, error)
	RemoveAll(name string) error
	Exists(path string) bool
	MkdirAll(path string, perm os.FileMode) error
}

type OsFileManager struct{}

func (OsFileManager) Create(name string) (File, error) {
	file, _ := os.Create(name)
	return &OsFile{file}, nil
}

func (OsFileManager) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

func (OsFileManager) Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (OsFileManager) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
