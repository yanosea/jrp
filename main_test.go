package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"
	"github.com/yanosea/jrp/test"
)

func TestMain(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()
	defaultDBFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")

	type want struct {
		exitCode int
		stdOut   string
		errOut   string
	}
	tests := []struct {
		name  string
		want  want
		setup func()
	}{
		{
			name: "positive testing (with no db file)",
			want: want{exitCode: 0, stdOut: constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED + "\n", errOut: ""},
			setup: func() {
				os.RemoveAll(defaultDBFileDirPath)
			},
		},
		{
			name: "positive testing (with db file)",
			want: want{exitCode: 0, stdOut: "\n", errOut: ""},
			setup: func() {
				tdl := logic.NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tdl.Download()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
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
	os.RemoveAll(defaultDBFileDirPath)
}
