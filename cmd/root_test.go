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
	"github.com/yanosea/jrp/internal/database"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/spinnerservice"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"

	mock_cmdwrapper "github.com/yanosea/jrp/mock/cmdwrapper"
	mock_generator "github.com/yanosea/jrp/mock/generator"
)

func TestNewGlobalOption(t *testing.T) {
	type args struct {
		out    io.Writer
		errOut io.Writer
		args   []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{out: os.Stdout, errOut: os.Stderr, args: []string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := cmd.NewGlobalOption(tt.args.out, tt.args.errOut, tt.args.args)
			if u == nil {
				t.Errorf("NewGlobalOption() : returned nil")
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()
	dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	tdl := logic.NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, spinnerservice.NewRealSpinnerService())

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
				if err := tdl.Download(); err != nil {
					t.Error(err)
				}
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
				if err := tdl.Download(); err != nil {
					t.Error(err)
				}
			},
		}, {
			name:    "negative testing (Execute() fails)",
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
				t.Errorf("Execute() : exit code = %v, want = %v", got, tt.want)
			}
		})
	}
	os.RemoveAll(dbFileDirPath)
}

func TestNewRootCommand(t *testing.T) {
	type args struct {
		out    io.Writer
		errOut io.Writer
		args   []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{out: os.Stdout, errOut: os.Stderr, args: []string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.NewRootCommand(tt.args.out, tt.args.errOut, tt.args.args)
			if got == nil {
				t.Errorf("NewRootCommand() : returned nil")
			}
		})
	}
}

func TestRootRunE(t *testing.T) {
	type args struct {
		o cmd.RootOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{o: cmd.RootOption{Out: os.Stdout, ErrOut: os.Stderr, Args: nil, Number: 0, Generator: logic.NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := tt.args.o.RootRunE(nil, nil)
			if err != nil {
				t.Errorf("RootRunE() : error = %v", err)
			}
		})
	}
}

func TestRootGenerate(t *testing.T) {
	type args struct {
		o cmd.RootOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name:    "positive testing",
			args:    args{o: cmd.RootOption{Out: os.Stdout, ErrOut: os.Stderr, Args: []string{""}, Number: 1, Generator: logic.NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})}},
			wantErr: false,
			setup:   nil,
		}, {
			name:    "negative testing (Generate() fails)",
			args:    args{o: cmd.RootOption{Out: os.Stdout, ErrOut: os.Stderr, Args: []string{""}, Number: 1, Generator: nil}},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mg := mock_generator.NewMockGenerator(mockCtrl)
				mg.EXPECT().Generate(tt.o.Number).Return(nil, errors.New("failed to generate japanese random phrase"))
				tt.o.Generator = mg
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.setup != nil {
				tt.setup(ctrl, &tt.args)
			}

			err := tt.args.o.RootGenerate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RootGenerate() : error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
