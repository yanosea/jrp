package testutility

import (
	"testing"

	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/os"
)

// Capturable is an interface for capturing output.
type Capturable interface {
	CaptureOutput() (string, string, error)
}

// Capturer is a struct for capturing output.
type Capturer struct {
	OutBuffer bufferproxy.Buffer
	ErrBuffer bufferproxy.Buffer
	OsProxy   osproxy.Os
}

// NewCapturer is a constructor for Capturer.
func NewCapturer(
	outBuffer bufferproxy.Buffer,
	errBuffer bufferproxy.Buffer,
	osProxy osproxy.Os,
) *Capturer {
	return &Capturer{
		OutBuffer: outBuffer,
		ErrBuffer: errBuffer,
		OsProxy:   osProxy,
	}
}

// CaptureOutput captures the output of the function and returns the captured.
func (c *Capturer) CaptureOutput(t *testing.T, fnc func()) (string, string, error) {
	t.Helper()

	// save the original stdout and stderr
	origStdout := osproxy.Stdout
	origStderr := osproxy.Stderr
	defer func() {
		osproxy.Stdout = origStdout
		osproxy.Stderr = origStderr
	}()

	// create a pipe for stdout and stderr
	rOut, wOut, _ := c.OsProxy.Pipe()
	rErr, wErr, _ := c.OsProxy.Pipe()
	osproxy.Stdout = wOut.FieldFile
	osproxy.Stderr = wErr.FieldFile

	// execute the function
	fnc()

	// close the pipe
	wOut.Close()
	wErr.Close()

	// read from the pipe of stdout
	if _, err := c.OutBuffer.ReadFrom(rOut); err != nil {
		return "", "", err
	}

	// read from the pipe of stderr
	if _, err := c.ErrBuffer.ReadFrom(rErr); err != nil {
		return "", "", err
	}

	// return the captured output
	stdout := c.OutBuffer.String()
	errout := c.ErrBuffer.String()

	// reset the buffer
	c.OutBuffer.Reset()
	c.ErrBuffer.Reset()

	return stdout, errout, nil
}
