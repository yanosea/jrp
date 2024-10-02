package osproxy

import (
	"os"

	"github.com/yanosea/jrp/app/proxy/fs"
)

// Os is an interface for os.
type Os interface {
	Create(name string) (FileInstanceInterface, error)
	FileMode(perm fsproxy.FileMode) fsproxy.FileMode
	Getenv(key string) string
	IsNotExist(err error) bool
	MkdirAll(path string, perm fsproxy.FileMode) error
	Open(name string) (*FileInstance, error)
	Pipe() (*FileInstance, *FileInstance, error)
	Remove(name string) error
	RemoveAll(path string) error
	Stat(name string) (*fsproxy.FileInfoInstance, error)
	TempDir() string
}

// OsProxy is a struct that implements Os.
type OsProxy struct{}

// New is a constructor for OsProxy.
func New() Os {
	return &OsProxy{}
}

// Create is a proxy for os.Create.
func (*OsProxy) Create(name string) (FileInstanceInterface, error) {
	file, _ := os.Create(name)
	return &FileInstance{FieldFile: file}, nil
}

// Filemode is a proxy for os.FileMode.
func (*OsProxy) FileMode(perm fsproxy.FileMode) fsproxy.FileMode {
	return perm
}

// Getenv is a proxy for os.Getenv.
func (*OsProxy) Getenv(key string) string {
	return os.Getenv(key)
}

// IsNotExist is a proxy for os.IsNotExist.
func (*OsProxy) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// MkdirAll is a proxy for os.MkdirAll.
func (*OsProxy) MkdirAll(path string, perm fsproxy.FileMode) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// Open is a proxy for os.Open.
func (*OsProxy) Open(name string) (*FileInstance, error) {
	file, _ := os.Open(name)
	return &FileInstance{FieldFile: file}, nil
}

// Pipe is a proxy for os.Pipe.
func (*OsProxy) Pipe() (*FileInstance, *FileInstance, error) {
	r, w, err := os.Pipe()
	return &FileInstance{FieldFile: r}, &FileInstance{FieldFile: w}, err
}

// Remove is a proxy for os.Remove.
func (*OsProxy) Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll is a proxy for os.RemoveAll.
func (*OsProxy) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Stat is a proxy for os.Stat.
func (*OsProxy) Stat(name string) (*fsproxy.FileInfoInstance, error) {
	fileInfo, err := os.Stat(name)
	return &fsproxy.FileInfoInstance{FieldFileInfo: fileInfo}, err
}

// TempDir is a proxy for os.TempDir.
func (*OsProxy) TempDir() string {
	return os.TempDir()
}
