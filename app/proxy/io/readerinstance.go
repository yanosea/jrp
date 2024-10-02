package ioproxy

import (
	"io"
)

// ReaderInstanceInterface is an interface for io.Reader.
type ReaderInstanceInterface interface {
	io.Reader
}

// ReaderInstance is a struct that implements ReaderInstanceInterface.
type ReaderInstance struct{}
