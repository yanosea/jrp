package proxy

import (
	"io"
)

// Io is an interface that provides a proxy of the methods of io.
type Io interface {
	Copy(dst io.Writer, src io.Reader) (int64, error)
}

// ioProxy is a proxy struct that implements the Io interface.
type ioProxy struct{}

// NewIo returns a new instance of the Io interface.
func NewIo() Io {
	return &ioProxy{}
}

// Copy copies from src to dst until either EOF is reached on src or an error occurs.
func (ioProxy) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
