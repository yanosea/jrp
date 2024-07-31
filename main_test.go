package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/test"
)

func TestMain(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()
	dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	os.RemoveAll(dbFileDirPath)

	type want struct {
		exitCode int
		stdOut   string
		errOut   string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive testing",
			want: want{exitCode: 0, stdOut: constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED + "\n", errOut: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origOsExit := osExit
			osExit = func(code int) {
				if code != tt.want.exitCode {
					t.Errorf("main() : exit code = %v, want = %v", code, tt.want.exitCode)
				}
			}
			defer func() {
				osExit = origOsExit
			}()

			stdOut, errOut := test.CaptureOutput(t, func() {
				main()
			})

			if stdOut != tt.want.stdOut {
				t.Errorf("main() :\n output = '%v', want = '%v'", stdOut, tt.want.stdOut)
			}
			if errOut != tt.want.errOut {
				t.Errorf("main() :\n error output = '%v', want = '%v'", errOut, tt.want.errOut)
			}
		})
	}
}
