package test

import (
	"bytes"
	"os"
	"testing"
)

func CaptureOutput(t *testing.T, fnc func()) (string, string) {
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

	var outBuf, errBuf bytes.Buffer

	outBuf.ReadFrom(rOut)
	errBuf.ReadFrom(rErr)

	stdout := outBuf.String()
	stderr := errBuf.String()

	return stdout, stderr
}
