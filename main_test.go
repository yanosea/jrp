package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"
)

func TestMain(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()
	dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	tdl := logic.NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})

	tests := []struct {
		name  string
		want  int
		setup func()
	}{
		{
			name: "positive testing (with no db file)",
			want: 0,
			setup: func() {
				os.RemoveAll(dbFileDirPath)
			},
		}, {
			name: "positive testing (with db file)",
			want: 0,
			setup: func() {
				os.RemoveAll(dbFileDirPath)
				if err := tdl.Download(); err != nil {
					t.Error(err)
				}
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
				if code != tt.want {
					t.Errorf("main() : exit code = %v, want = %v", code, tt.want)
				}
			}
			defer func() {
				osExit = origOsExit
			}()

			main()
		})
	}
}
