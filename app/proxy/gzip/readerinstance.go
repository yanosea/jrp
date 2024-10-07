package gzipproxy

import (
	"compress/gzip"
	"github.com/yanosea/jrp/app/proxy/io"
)

// ReaderInstanceInterface is an interface for gzip.Reader.
type ReaderInstanceInterface interface {
	ioproxy.ReaderInstanceInterface
	Close() error
}

// ReaderInstance is a struct that implements ReaderInstanceInterface.
type ReaderInstance struct {
	FieldReader gzip.Reader
}
