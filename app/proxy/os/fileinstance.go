package osproxy

import (
	"os"
)

// FileInstanceInterface is an interface for os.File.
type FileInstanceInterface interface {
	Close() error
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Write(b []byte) (n int, err error)
}

// FileInstance is a struct that implements FileInstanceInterface.
type FileInstance struct {
	FieldFile *os.File
}

// Close is a proxy for os.File.Close().
func (f *FileInstance) Close() error {
	return f.FieldFile.Close()
}

// Read is a proxy for os.File.Read().
func (f *FileInstance) Read(p []byte) (n int, err error) {
	return f.FieldFile.Read(p)
}

// Seek is a proxy for os.File.Seek().
func (f *FileInstance) Seek(offset int64, whence int) (int64, error) {
	return f.FieldFile.Seek(offset, whence)
}

// Write is a proxy for os.File.Write().
func (f *FileInstance) Write(b []byte) (n int, err error) {
	return f.FieldFile.Write(b)
}
