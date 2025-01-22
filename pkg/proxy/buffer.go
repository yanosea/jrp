package proxy

import (
	"bytes"
	"io"
)

// Buffer is an interface that provides a proxy of the methods of bytes.Buffer.
type Buffer interface {
	ReadFrom(r io.Reader) (int64, error)
	Reset()
	String() string
}

// bufferProxy is a proxy struct that implements the Buffer interface.
type bufferProxy struct {
	bytes.Buffer
}

// NewBuffer returns a new instance of the Buffer interface.
func NewBuffer() Buffer {
	return &bufferProxy{}
}

// ReadFrom is a proxy method that calls the ReadFrom method of the bytes.Buffer.
func (b *bufferProxy) ReadFrom(r io.Reader) (int64, error) {
	return b.Buffer.ReadFrom(r)
}

// Reset is a proxy method that calls the Reset method of the bytes.Buffer.
func (b *bufferProxy) Reset() {
	b.Buffer.Reset()
}

// String is a proxy method that calls the String method of the bytes.Buffer.
func (b *bufferProxy) String() string {
	return b.Buffer.String()
}
