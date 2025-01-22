package proxy

import (
	"os"
)

// Os is an interface that provides a proxy of the methods of os.
type Os interface {
	Create(name string) (File, error)
	Getenv(key string) string
	IsNotExist(err error) bool
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (File, error)
	RemoveAll(path string) error
	Rename(oldpath, newpath string) error
	Stat(name string) (os.FileInfo, error)
	TempDir() string
	UserHomeDir() (string, error)
}

// osProxy is a proxy struct that implements the Os interface.
type osProxy struct{}

// NewOs returns a new instance of the Os interface.
func NewOs() Os {
	return &osProxy{}
}

// Create creates or truncates the named file.
func (osProxy) Create(name string) (File, error) {
	file, err := os.Create(name)
	return &fileProxy{file}, err
}

// Getenv returns the value of the environment variable named by the key.
func (osProxy) Getenv(key string) string {
	return os.Getenv(key)
}

// IsNotExist reports whether the error is known to report that a file or directory does not exist.
func (osProxy) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
func (osProxy) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Open opens the named file for reading.
func (osProxy) Open(name string) (File, error) {
	file, err := os.Open(name)
	return &fileProxy{file}, err
}

// RemoveAll removes path and any children it contains.
func (osProxy) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Rename renames (moves) oldpath to newpath.
func (osProxy) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// Stat returns a FileInfo describing the named file.
func (osProxy) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// TempDir returns the default directory to use for temporary files.
func (osProxy) TempDir() string {
	return os.TempDir()
}

// UserHomeDir returns the current user's home directory.
func (osProxy) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

// File is an interface that provides a proxy of the methods of os.File.
type File interface {
	Close() error
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

// fileProxy is a proxy struct that implements the File interface.
type fileProxy struct {
	file *os.File
}

// Close closes the File, rendering it unusable for I/O.
func (f *fileProxy) Close() error {
	return f.file.Close()
}

// Read reads up to len(p) bytes into p.
func (f *fileProxy) Read(b []byte) (n int, err error) {
	return f.file.Read(b)
}

// Write writes len(p) bytes from p to the underlying data stream.
func (f *fileProxy) Write(b []byte) (n int, err error) {
	return f.file.Write(b)
}
