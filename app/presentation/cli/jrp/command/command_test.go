package command

import (
	"context"
	"errors"
	o "os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/fatih/color"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
	"go.uber.org/mock/gomock"
)

func Test_newCli(t *testing.T) {
	cobra := proxy.NewCobra()

	type args struct {
		cobra proxy.Cobra
	}
	tests := []struct {
		name string
		args args
		want *cli
	}{
		{
			name: "positive testing",
			args: args{
				cobra: cobra,
			},
			want: &cli{
				Cobra:       cobra,
				RootCommand: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCli(tt.args.cobra); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCli() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cli_Init(t *testing.T) {
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	cobra := proxy.NewCobra()
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(o.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type fields struct {
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func(mockCtrl *gomock.Controller)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(_ *gomock.Controller) {
					c := &cli{
						Cobra:             cobra,
						RootCommand:       nil,
						ConnectionManager: nil,
					}
					if got := c.Init(
						proxy.NewEnvconfig(),
						proxy.NewSql(),
						"0.0.0",
						utility.NewFileUtil(
							proxy.NewGzip(),
							proxy.NewIo(),
							proxy.NewOs(),
						),
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 0 {
						t.Errorf("cli.Init() = %v, want %v", got, 0)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := o.Setenv("JRP_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_DB", filepath.Join(o.TempDir(), "jrp.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
			},
			cleanup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
		{
			name: "negative testing (configurator.GetConfig() failed)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
					mockEnvconfig.EXPECT().Process("", gomock.Any()).Return(errors.New("EnvconfigProxy.Process() failed"))
					c := &cli{
						Cobra:             cobra,
						RootCommand:       nil,
						ConnectionManager: nil,
					}
					if got := c.Init(
						mockEnvconfig,
						proxy.NewSql(),
						"0.0.0",
						utility.NewFileUtil(
							proxy.NewGzip(),
							proxy.NewIo(),
							proxy.NewOs(),
						),
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 1 {
						t.Errorf("cli.Init() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : EnvconfigProxy.Process() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := o.Setenv("JRP_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_DB", filepath.Join(o.TempDir(), "jrp.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
			},
			cleanup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
		{
			name: "negative testing (JrpDB InitializeConnection failed)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockConnectionManager := database.NewMockConnectionManager(mockCtrl)
					mockConnectionManager.EXPECT().InitializeConnection(gomock.Any()).Return(errors.New("ConnectionManager.InitializeConnection() failed"))
					c := &cli{
						Cobra:             cobra,
						RootCommand:       nil,
						ConnectionManager: mockConnectionManager,
					}
					if got := c.Init(
						proxy.NewEnvconfig(),
						proxy.NewSql(),
						"0.0.0",
						utility.NewFileUtil(
							proxy.NewGzip(),
							proxy.NewIo(),
							proxy.NewOs(),
						),
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 1 {
						t.Errorf("cli.Init() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : ConnectionManager.InitializeConnection() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := o.Setenv("JRP_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_DB", filepath.Join(o.TempDir(), "jrp.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
			},
			cleanup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
		{
			name: "negative testing (WNJpnDB InitializeConnection failed)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockConnectionManager := database.NewMockConnectionManager(mockCtrl)
					mockConnectionManager.EXPECT().InitializeConnection(gomock.Any()).Return(nil)
					mockConnectionManager.EXPECT().InitializeConnection(gomock.Any()).Return(errors.New("ConnectionManager.InitializeConnection() failed"))
					c := &cli{
						Cobra:             cobra,
						RootCommand:       nil,
						ConnectionManager: mockConnectionManager,
					}
					if got := c.Init(
						proxy.NewEnvconfig(),
						proxy.NewSql(),
						"0.0.0",
						utility.NewFileUtil(
							proxy.NewGzip(),
							proxy.NewIo(),
							proxy.NewOs(),
						),
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 1 {
						t.Errorf("cli.Init() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : ConnectionManager.InitializeConnection() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := o.Setenv("JRP_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_DB", filepath.Join(o.TempDir(), "jrp.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
			},
			cleanup: func() {
				output = ""
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			output = ""
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			c := utility.NewCapturer(tt.fields.StdBuffer, tt.fields.ErrBuffer)
			gotStdOut, gotStdErr, err := c.CaptureOutput(func() {
				tt.args.fnc(mockCtrl)
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Init() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Init() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}

func Test_cli_Run(t *testing.T) {
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()

	type fields struct {
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func(mockCtrl *gomock.Controller)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func()
		clear      func()
	}{
		{
			name: "positive testing",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockConnectionManager := database.NewMockConnectionManager(mockCtrl)
					mockConnectionManager.EXPECT().CloseAllConnections().Return(nil)
					c := &cli{
						Cobra:             proxy.NewCobra(),
						RootCommand:       mockCommand,
						ConnectionManager: mockConnectionManager,
					}
					if got := c.Run(context.Background()); got != 0 {
						t.Errorf("cli.Run() = %v, want %v", got, 0)
					}
				},
			},
			wantStdOut: "\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			clear: func() {
				output = ""
			},
		},
		{
			name: "negative testing (c.RootCommand.ExecuteContext() failed)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(errors.New("CommandProxy.ExecuteContext() failed"))
					mockConnectionManager := database.NewMockConnectionManager(mockCtrl)
					mockConnectionManager.EXPECT().CloseAllConnections().Return(nil)
					c := &cli{
						Cobra:             proxy.NewCobra(),
						RootCommand:       mockCommand,
						ConnectionManager: mockConnectionManager,
					}
					if got := c.Run(context.Background()); got != 1 {
						t.Errorf("cli.Run() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : CommandProxy.ExecuteContext() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			clear: func() {
				output = ""
			},
		},
		{
			name: "negative testing (c.ConnectionManager.CloseAllConnections() failed)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockConnectionManager := database.NewMockConnectionManager(mockCtrl)
					mockConnectionManager.EXPECT().CloseAllConnections().Return(errors.New("ConnectionManager.CloseAllConnections() failed"))
					c := &cli{
						Cobra:             proxy.NewCobra(),
						RootCommand:       mockCommand,
						ConnectionManager: mockConnectionManager,
					}
					if got := c.Run(context.Background()); got != 1 {
						t.Errorf("cli.Run() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "\n",
			wantStdErr: color.RedString("Error : ConnectionManager.CloseAllConnections() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			clear: func() {
				output = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			output = ""
			c := utility.NewCapturer(tt.fields.StdBuffer, tt.fields.ErrBuffer)
			gotStdOut, gotStdErr, err := c.CaptureOutput(func() {
				tt.args.fnc(mockCtrl)
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Run() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Run() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}
