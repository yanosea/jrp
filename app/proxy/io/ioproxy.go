package ioproxy

import (
	"io"
)

// Io is an interface for io.
type Io interface {
	Copy(dst WriterInstanceInterface, src ReaderInstanceInterface) (int64, error)
}

// IoProxy is a struct that implements Io.
type IoProxy struct{}

// New is a constructor for IoProxy.
func New() Io {
	return &IoProxy{}
}

// Copy is a proxy for io.Copy.
func (*IoProxy) Copy(dst WriterInstanceInterface, src ReaderInstanceInterface) (int64, error) {
	return io.Copy(dst, src)
}
