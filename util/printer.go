package util

import (
	"fmt"
	"io"
)

func PrintlnWithWriter(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("%s", a[0]))
}
