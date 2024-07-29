package logic

import (
	"compress/gzip"
	"io"
)

type Gzip interface {
	NewReader(r io.Reader) (io.ReadCloser, error)
}

type DefaultGzip struct{}

func (DefaultGzip) NewReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}
