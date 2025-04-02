package utility

import (
	"os"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// Capturer is an interface that captures the output of a function.
type Capturer interface {
	CaptureOutput(fnc func()) (string, string, error)
}

// capturer is a struct that implements the Captures interface.
type capturer struct {
	// os is an interface for operating system functions.
	os proxy.Os
	// stdBuffer is a buffer for standard output.
	stdBuffer proxy.Buffer
	// errBuffer is a buffer for error output.
	errBuffer proxy.Buffer
}

// NewCapturer returns a new instance of the capturer struct.
func NewCapturer(
	os proxy.Os,
	stdBuffer proxy.Buffer,
	errBuffer proxy.Buffer,
) *capturer {
	return &capturer{
		os:        os,
		stdBuffer: stdBuffer,
		errBuffer: errBuffer,
	}
}

// CaptureOutput captures the output of a function.
func (c *capturer) CaptureOutput(fnc func()) (string, string, error) {
	origStdout := os.Stdout
	origStderr := os.Stderr
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	rOut, wOut, _ := c.os.Pipe()
	rErr, wErr, _ := c.os.Pipe()
	os.Stdout = wOut.(interface{ AsOsFile() *os.File }).AsOsFile()
	os.Stderr = wErr.(interface{ AsOsFile() *os.File }).AsOsFile()

	fnc()

	if err := wOut.Close(); err != nil {
		return "", "", err
	}
	if err := wErr.Close(); err != nil {
		return "", "", err
	}

	if _, err := c.stdBuffer.ReadFrom(rOut); err != nil {
		return "", "", err
	}

	if _, err := c.errBuffer.ReadFrom(rErr); err != nil {
		return "", "", err
	}

	stdout := c.stdBuffer.String()
	errout := c.errBuffer.String()

	c.stdBuffer.Reset()
	c.errBuffer.Reset()

	return stdout, errout, nil
}
