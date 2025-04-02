package presenter

import (
	"fmt"
	"io"
)

// Print is a function that writes the output to the writer.
func Print(writer io.Writer, output string) error {
	if output != "" && output != "\n" {
		_, err := fmt.Fprintf(writer, "%s\n", output)
		return err
	} else {
		_, err := fmt.Fprintln(writer)
		return err
	}
}
