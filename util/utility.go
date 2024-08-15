package util

import (
	"fmt"
	"io"
)

func PrintlnWithWriter(w io.Writer, a ...any) {
	fmt.Fprintf(w, fmt.Sprintf("%s", a[0])+"\n")
}
