package buffer

import (
	"bytes"
	"os"
)

type Buffer interface {
	ReadFrom(f *os.File) (int64, error)
	String() string
}

type DefaultBuffer struct {
	bytes.Buffer
}

func (b *DefaultBuffer) ReadFrom(f *os.File) (int64, error) {
	return b.Buffer.ReadFrom(f)
}

func (b *DefaultBuffer) String() string {
	return b.Buffer.String()
}
