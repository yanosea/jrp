package iomanager

import (
	"io"
)

type IOHelper interface {
	Copy(dst io.Writer, src io.Reader) (int64, error)
}

type DefaultIOHelper struct{}

func (DefaultIOHelper) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
