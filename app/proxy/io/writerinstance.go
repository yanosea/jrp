package ioproxy

import (
	"io"
)

// WriterInstanceInterface is an interface for io.Writer.
type WriterInstanceInterface interface {
	io.Writer
}

// WriterInstance is a struct that implements WriterInterface.
type WriterInstance struct{}
