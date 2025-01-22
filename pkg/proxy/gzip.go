package proxy

import (
	"compress/gzip"
	"io"
)

type Gzip interface {
	NewReader(r io.Reader) (GzipReader, error)
}

type gzipProxy struct{}

func NewGzip() Gzip {
	return &gzipProxy{}
}

func (gzipProxy) NewReader(r io.Reader) (GzipReader, error) {
	reader, err := gzip.NewReader(r)
	return &gzipReaderProxy{
		gzipReader: reader,
	}, err
}

// GzipReader is an interface that contains the utility functions for reading gzipped files.
type GzipReader interface {
	Close() error
	Read(p []byte) (n int, err error)
}

// gzipReaderProxy is a struct that contains the utility functions for reading gzipped files.
type gzipReaderProxy struct {
	gzipReader *gzip.Reader
}

// Close closes the gzip.Reader.
func (g *gzipReaderProxy) Close() error {
	return g.gzipReader.Close()
}

// Read reads up to len(p) bytes into p.
func (g *gzipReaderProxy) Read(p []byte) (n int, err error) {
	return g.gzipReader.Read(p)
}
