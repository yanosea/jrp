package fs

import (
	"os"
)

type File interface {
	Seek(offset int64, whence int) (int64, error)
	Close() error
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
}

type OsFile struct {
	*os.File
}

func (f *OsFile) Seek(offset int64, whence int) (int64, error) {
	return f.File.Seek(offset, whence)
}

func (f *OsFile) Close() error {
	return f.File.Close()
}
