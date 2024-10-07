package gzipproxy

import (
	"compress/gzip"

	"github.com/yanosea/jrp/app/proxy/io"
)

// Gzip is an interface for gzip.
type Gzip interface {
	NewReader(r ioproxy.ReaderInstanceInterface) (ReaderInstanceInterface, error)
}

// GzipProxy is a struct that implements Gzip.
type GzipProxy struct{}

// New is a constructor for GzipProxy.
func New() Gzip {
	return &GzipProxy{}
}

// NewReader is a proxy for gzip.NewReader.
func (*GzipProxy) NewReader(r ioproxy.ReaderInstanceInterface) (ReaderInstanceInterface, error) {
	return gzip.NewReader(r)
}
