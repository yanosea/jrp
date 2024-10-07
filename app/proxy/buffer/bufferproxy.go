package bufferproxy

import (
	"bytes"

	"github.com/yanosea/jrp/app/proxy/io"
)

// Buffer is an interface for buffer.
type Buffer interface {
	ReadFrom(r ioproxy.ReaderInstanceInterface) (int64, error)
	Reset()
	String() string
	Write(p []byte) (n int, err error)
}

// BufferProxy is a struct that implements Buffer.
type BufferProxy struct {
	bytes.Buffer
}

// New is a constructor for BufferProxy.
func New() Buffer {
	return &BufferProxy{}
}

// ReadFrom is a proxy for buffer.ReadFrom.
func (b *BufferProxy) ReadFrom(r ioproxy.ReaderInstanceInterface) (int64, error) {
	return b.Buffer.ReadFrom(r)
}

// Reset is a proxy for buffer.Reset.
func (b *BufferProxy) Reset() {
	b.Buffer.Reset()
}

// String is a proxy for buffer.String.
func (b *BufferProxy) String() string {
	return b.Buffer.String()
}

// Write is a proxy for buffer.Write.
func (b *BufferProxy) Write(p []byte) (n int, err error) {
	return b.Buffer.Write(p)
}
