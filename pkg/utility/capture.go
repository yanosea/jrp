package utility

import (
	"os"

	"github.com/yanosea/jrp/pkg/proxy"
)

// Capturer is an interface that captures the output of a function.
type Capturer interface {
	CaptureOutput(fnc func()) (string, string, error)
}

// capturer is a struct that implements the Captures interface.
type capturer struct {
	// stdBuffer is a buffer for standard output.
	stdBuffer proxy.Buffer
	// errBuffer is a buffer for error output.
	errBuffer proxy.Buffer
}

// NewCapturer returns a new instance of the capturer struct.
func NewCapturer(
	stdBuffer proxy.Buffer,
	errBuffer proxy.Buffer,
) *capturer {
	return &capturer{
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

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	fnc()

	wOut.Close()
	wErr.Close()

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
