package util

import (
	"fmt"
	"io"
)

func FormatIndent(m string) string {
	return "  " + m
}

func PrintlnWithWriter(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("%s", a[0]))
}

func PrintWithWriterWithBlankLineBelow(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("%s\n", a[0]))
}

func PrintWithWriterWithBlankLineAbove(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("\n%s", a[0]))
}

func PrintWithWriterBetweenBlankLine(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("\n%s\n", a[0]))
}
