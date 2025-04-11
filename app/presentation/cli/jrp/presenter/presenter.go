package presenter

import (
	"fmt"
	"io"
)

// PrintFunc is a function type that defines the signature for printing output.
type PrintFunc func(writer io.Writer, output string) error

// Print is a function that writes the output to the writer.
var Print PrintFunc = func(writer io.Writer, output string) error {
	if output != "" && output != "\n" {
		_, err := fmt.Fprintf(writer, "%s\n", output)
		return err
	} else {
		_, err := fmt.Fprintln(writer)
		return err
	}
}
