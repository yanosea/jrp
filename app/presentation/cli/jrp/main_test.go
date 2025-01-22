package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"

	"github.com/yanosea/jrp/app/presentation/cli/jrp/command"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"

	"go.uber.org/mock/gomock"
)

func Test_main(t *testing.T) {
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	os.Setenv("JRP_DB_TYPE", "sqlite")
	os.Setenv("JRP_DB", filepath.Join(os.TempDir(), "jrp.db"))
	os.Setenv("JRP_WNJPN_DB_TYPE", "sqlite")
	os.Setenv("JRP_WNJPN_DB", filepath.Join(os.TempDir(), "wnjpn.db"))
	origExit := exit
	exit = func(code int) {}
	defer func() {
		exit = origExit
	}()
	origArgs := os.Args

	type fields struct {
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
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp"}
					defer func() {
						os.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: color.YellowString("⚡ You have to execute \"download\" to use jrp...") + "\n",
			wantStdErr: "",
			setup: func(_ *gomock.Controller) {
				if err := os.Remove(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
			cleanup: nil,
		},
		{
			name: "download execution (jrp download)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "download"}
					defer func() {
						os.Args = origArgs
					}()
					main()
				},
			},
			wantStdOut: color.GreenString("✅ Downloaded successfully! Now, you are ready to use jrp!") + "\n",
			wantStdErr: "",
			setup: func(_ *gomock.Controller) {
				if err := os.Remove(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
			cleanup: nil,
		},
		{
			name: "re:download execution (jrp download)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "download"}
					defer func() {
						os.Args = origArgs
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
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "generate"}
					defer func() {
						os.Args = origArgs
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
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "generate", "10"}
					defer func() {
						os.Args = origArgs
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
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "history", "--all"}
					defer func() {
						os.Args = origArgs
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
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "favorite", "1"}
					defer func() {
						os.Args = origArgs
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
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "unfavorite", "1"}
					defer func() {
						os.Args = origArgs
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
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp", "history", "clear", "--force", "--no-confirm"}
					defer func() {
						os.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: color.GreenString("✅ Cleared successfully!") + "\n",
			wantStdErr: "",
			setup:      nil,
			cleanup: func() {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (cli.Init() failed)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					os.Args = []string{"/path/to/jrp"}
					defer func() {
						os.Args = origArgs
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
			c := utility.NewCapturer(tt.fields.StdBuffer, tt.fields.ErrBuffer)
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
