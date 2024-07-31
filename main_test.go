package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/usermanager"
)

func TestMain(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()

	type args struct {
		command         string
		commandLineArgs []string
	}
	type want struct {
		exitCode int
		stdOut   string
		errOut   string
	}
	tests := []struct {
		name  string
		args  args
		want  want
		setup func(tt *args)
	}{
		{
			name: "positive testing (no command, no command line args, not downloaded the db file yet)",
			args: args{command: "", commandLineArgs: []string{""}},
			want: want{exitCode: 0, stdOut: constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED + "\n", errOut: ""},
			setup: func(tt *args) {
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(&tt.args)
			}

			origOsExit := osExit
			osExit = func(code int) {
				if code != tt.want.exitCode {
					t.Fatalf("osExit was called with code %v", code)
					t.Errorf("main() exit code = %v, want = %v", code, tt.want.exitCode)
				}
			}
			defer func() {
				osExit = origOsExit
			}()

			cla := append([]string{tt.args.command}, tt.args.commandLineArgs...)
			os.Args = cla
			stdOut, errOut := CaptureOutput(t, main)
			if stdOut != tt.want.stdOut {
				t.Errorf("main() output '%v', want '%v'", stdOut, tt.want.stdOut)
			}
			if errOut != tt.want.errOut {
				t.Errorf("main() error output '%v', want '%v'", errOut, tt.want.errOut)
			}
		})
	}
}

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

	if _, err := outBuf.ReadFrom(rOut); err != nil {
		t.Fatalf("fail read stdout: %v", err)
	}
	if _, err := errBuf.ReadFrom(rErr); err != nil {
		t.Fatalf("fail read stderr: %v", err)
	}

	stdout := outBuf.String()
	stderr := errBuf.String()

	return stdout, stderr
}
