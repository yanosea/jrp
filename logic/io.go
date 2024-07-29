package logic

import (
	"io"
)

type IO interface {
	Copy(dst io.Writer, src io.Reader) (int64, error)
}

type DefaultIO struct{}

func (DefaultIO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
