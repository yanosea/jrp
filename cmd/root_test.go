package cmd_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/cmd"
	"github.com/yanosea/jrp/internal/cmdwrapper"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"

	mock_cmdwrapper "github.com/yanosea/jrp/mock/cmdwrapper"
)

func TestExecute(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()
	dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	tdl := logic.NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})

	type args struct {
		globalOption *cmd.GlobalOption
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
		setup   func(mockCmd *mock_cmdwrapper.MockICommand)
	}{
		{
			name:    "positive testing (with no args, no db file)",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{})},
			want:    0,
			wantErr: false,
			setup: func(_ *mock_cmdwrapper.MockICommand) {
				os.RemoveAll(dbFileDirPath)
			},
		}, {
			name:    "positive testing (with no args, with db file)",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{})},
			want:    0,
			wantErr: false,
			setup: func(_ *mock_cmdwrapper.MockICommand) {
				os.RemoveAll(dbFileDirPath)
				tdl.Download()
			},
		}, {
			name:    "positive testing (with args, no db file)",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{"2"})},
			want:    0,
			wantErr: false,
			setup: func(_ *mock_cmdwrapper.MockICommand) {
				os.RemoveAll(dbFileDirPath)
			},
		}, {
			name:    "positive testing (with args, db file)",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{"2"})},
			want:    0,
			wantErr: false,
			setup: func(_ *mock_cmdwrapper.MockICommand) {
				os.RemoveAll(dbFileDirPath)
				tdl.Download()
			},
		}, {
			name:    "negative testing (rootCmd.Execute() fails)",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{})},
			want:    1,
			wantErr: true,
			setup: func(mockCmd *mock_cmdwrapper.MockICommand) {
				mockCmd.EXPECT().Execute().Return(errors.New("failed to execute command"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.setup != nil && !tt.wantErr {
				tt.setup(nil)
			}

			if tt.setup != nil && tt.wantErr {
				mockCmd := mock_cmdwrapper.NewMockICommand(ctrl)
				tt.setup(mockCmd)
				tt.args.globalOption.NewRootCommand = func(ow, ew io.Writer, cmdArgs []string) cmdwrapper.ICommand {
					return mockCmd
				}
			}

			if got := tt.args.globalOption.Execute(); (got != 0) != tt.wantErr {
				t.Errorf("Execute() = %v, want = %v", got, tt.want)
			}
		})
	}
}
