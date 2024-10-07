package filepathproxy

import (
	"path/filepath"
)

// FilePath is an interface for filepath.
type FilePath interface {
	Join(elem ...string) string
}

// FilePathProxy is a struct that implements FilePath.
type FilePathProxy struct{}

// New is a constructor for FilepathProxy
func New() FilePath {
	return &FilePathProxy{}
}

// Join is a proxy for filepath.Join.
func (*FilePathProxy) Join(elem ...string) string {
	return filepath.Join(elem...)
}
