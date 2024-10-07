package fsproxy

import (
	"io/fs"
)

// FileInfoInstanceInterface is an interface for fs.FileInfo.
type FileInfoInstanceInterface interface{}

// FileInfoInstance is a struct that implements FileInfoInstanceInterface.
type FileInfoInstance struct {
	FieldFileInfo fs.FileInfo
}

// FileMode is for fs.FileMode.
type FileMode fs.FileMode
