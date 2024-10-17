package filepathproxy

import (
	"path/filepath"
)

// FilePath is an interface for filepath.
type FilePath interface {
	Base(path string) string
	Dir(path string) string
	Join(elem ...string) string
}

// FilePathProxy is a struct that implements FilePath.
type FilePathProxy struct{}

// New is a constructor for FilepathProxy
func New() FilePath {
	return &FilePathProxy{}
}

func (*FilePathProxy) Base(path string) string {
	return filepath.Base(path)
}

// Dir is a proxy for filepath.Dir.
func (*FilePathProxy) Dir(path string) string {
	return filepath.Dir(path)
}

// Join is a proxy for filepath.Join.
func (*FilePathProxy) Join(elem ...string) string {
	return filepath.Join(elem...)
}
