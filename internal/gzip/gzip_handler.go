package gzip

import (
	"compress/gzip"
	"io"
)

type GzipHandler interface {
	NewReader(r io.Reader) (io.ReadCloser, error)
}

type DefaultGzipHandler struct{}

func (DefaultGzipHandler) NewReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}
