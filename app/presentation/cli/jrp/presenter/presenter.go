package presenter

import (
	"fmt"
	"io"
)

// Print is a function that writes the output to the writer.
func Print(writer io.Writer, output string) {
	if output != "" && output != "\n" {
		fmt.Fprintf(writer, "%s\n", output)
	} else {
		fmt.Fprintln(writer)
	}
}
