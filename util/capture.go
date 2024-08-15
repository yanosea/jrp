package util

import (
	"os"
	"testing"

	"github.com/yanosea/jrp/internal/buffer"
)

type Capturer interface {
	CaptureOutput() (string, string, error)
}

type DefaultCapturer struct {
	OutBuffer buffer.Buffer
	ErrBuffer buffer.Buffer
}

func NewCapturer(ob buffer.Buffer, eb buffer.Buffer) *DefaultCapturer {
	return &DefaultCapturer{
		OutBuffer: ob,
		ErrBuffer: eb,
	}
}

func (d *DefaultCapturer) CaptureOutput(t *testing.T, fnc func()) (string, string, error) {
	t.Helper()

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

	if _, err := d.OutBuffer.ReadFrom(rOut); err != nil {
		return "", "", err
	}

	if _, err := d.ErrBuffer.ReadFrom(rErr); err != nil {
		return "", "", err
	}

	stdout := d.OutBuffer.String()
	stderr := d.ErrBuffer.String()

	return stdout, stderr, nil
}
