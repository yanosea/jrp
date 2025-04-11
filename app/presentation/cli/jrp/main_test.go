package main

import (
	o "os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"

	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/command"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func Test_main(t *testing.T) {
	os := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	if err := o.Setenv("JRP_DB_TYPE", "sqlite"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := o.Setenv("JRP_DB", filepath.Join(o.TempDir(), "jrp.db")); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := o.Setenv("JRP_WNJPN_DB_TYPE", "sqlite"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := o.Setenv("JRP_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	origExit := exit
	exit = func(code int) {}
	defer func() {
		exit = origExit
	}()
	origArgs := o.Args

	type fields struct {
		Os        proxy.Os
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func()
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantAny    bool
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller)
		cleanup    func()
	}{
		{
			name: "initial execution (jrp)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: color.YellowString("⚡ You have to execute \"download\" to use jrp...") + "\n",
			wantStdErr: "",
			setup: func(_ *gomock.Controller) {
				if err := o.Remove(filepath.Join(o.TempDir(), "wnjpn.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
			cleanup: nil,
		},
		{
			name: "download execution (jrp download)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "download"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantStdOut: color.GreenString("✅ Downloaded successfully! Now, you are ready to use jrp!") + "\n",
			wantStdErr: "",
			setup: func(_ *gomock.Controller) {
				if err := o.Remove(filepath.Join(o.TempDir(), "wnjpn.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
			cleanup: nil,
		},
		{
			name: "re:download execution (jrp download)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "download"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: color.GreenString("✅ You are already ready to use jrp!") + "\n",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "generate execution (jrp generate)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "generate"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    true,
			wantStdOut: "",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "generate execution with arg 10 (jrp generate 10)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "generate", "10"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    true,
			wantStdOut: "",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "history execution with all options (jrp history --all)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "history", "--all"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    true,
			wantStdOut: "",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "favorite execution with with arg 1 (jrp favorite 1)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "favorite", "1"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: color.GreenString("✅ Favorited successfully!") + "\n",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "unfavorite execution with with arg 1 (jrp unfavorite 1)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "unfavorite", "1"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: color.GreenString("✅ Unfavorited successfully!") + "\n",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "clear execution with force option and no-confirm option (jrp history clear --force --no-confirm)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp", "history", "clear", "--force", "--no-confirm"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: color.GreenString("✅ Cleared successfully!") + "\n",
			wantStdErr: "",
			setup:      nil,
			cleanup: func() {
				if err := o.Remove(filepath.Join(o.TempDir(), "jrp.db")); err != nil && !o.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (cli.Init() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/jrp"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: "",
			wantStdErr: "",
			setup: func(mockCtrl *gomock.Controller) {
				mockCli := command.NewMockCli(mockCtrl)
				mockCli.EXPECT().Init(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1)
				mockCli.EXPECT().Run(gomock.Any())
				origNewCli := command.NewCli
				command.NewCli = func(cobra proxy.Cobra) command.Cli {
					return mockCli
				}
				t.Cleanup(func() {
					command.NewCli = origNewCli
				})
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			c := utility.NewCapturer(tt.fields.Os, tt.fields.StdBuffer, tt.fields.ErrBuffer)
			gotStdOut, gotStdErr, err := c.CaptureOutput(tt.args.fnc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantAny {
				if gotStdOut != tt.wantStdOut {
					t.Errorf("main() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
				}
				if gotStdErr != tt.wantStdErr {
					t.Errorf("main() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
				}
			} else {
				t.Logf("main() gotStdOut = %v", gotStdOut)
				t.Logf("main() gotStdErr = %v", gotStdErr)
			}
		})
	}
}
