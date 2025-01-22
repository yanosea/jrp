package command

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
	baseConfig "github.com/yanosea/jrp/app/config"
	"github.com/yanosea/jrp/app/infrastructure/database"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/command/jrp"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/command/jrp/generate"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/config"

	"github.com/yanosea/jrp/pkg/proxy"
)

func TestNewRootCommand(t *testing.T) {
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type args struct {
		cobra   proxy.Cobra
		version string
		conf    *config.JrpCliConfig
		output  *string
	}
	tests := []struct {
		name  string
		args  args
		setup func()
		clear func()
	}{
		{
			name: "positive testing",
			args: args{
				cobra:   proxy.NewCobra(),
				version: "0.0.0",
				conf: &config.JrpCliConfig{
					JrpConfig: baseConfig.JrpConfig{
						JrpDBType:   "sqlite",
						JrpDBDsn:    filepath.Join(os.TempDir(), "jrp.db"),
						WNJpnDBType: "sqlite",
						WNJpnDBDsn:  filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				},
				output: &output,
			},
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
			got := NewRootCommand(tt.args.cobra, tt.args.version, tt.args.conf, tt.args.output)
			if got == nil {
				t.Errorf("NewRootCommand() = %v, want not nil", got)
			} else {
				if err := got.RunE(nil, []string{}); err != nil {
					t.Errorf("Failed to run the root command: %v", err)
				}
			}
		})
	}
}

func Test_runRoot(t *testing.T) {
	type args struct {
		cmd         *c.Command
		cmdArgs     []string
		generateCmd proxy.Command
		versionCmd  proxy.Command
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setup   func(tt *args)
		cleanup func()
	}{
		{
			name: "positive testing (rootOps.Version is true)",
			args: args{
				cmd:     &c.Command{},
				cmdArgs: []string{},
				generateCmd: generate.NewGenerateCommand(
					proxy.NewCobra(),
					generate.NewInteractiveCommand(proxy.NewCobra(),
						&output,
					),
					&output,
				),
				versionCmd: jrp.NewVersionCommand(
					proxy.NewCobra(),
					"0.0.0",
					&output,
				),
			},
			want:    "jrp version 0.0.0",
			wantErr: false,
			setup: func(_ *args) {
				rootOps.Version = true
				output = ""
			},
			cleanup: func() {
				rootOps.Version = false
				output = ""
			},
		},
		{
			name: "positive testing (rootOps.Version is false)",
			args: args{
				cmd:     nil,
				cmdArgs: []string{},
				generateCmd: generate.NewGenerateCommand(
					proxy.NewCobra(),
					generate.NewInteractiveCommand(proxy.NewCobra(),
						&output,
					),
					&output,
				),
				versionCmd: jrp.NewVersionCommand(
					proxy.NewCobra(),
					"0.0.0",
					&output,
				),
			},
			want:    "",
			wantErr: false,
			setup: func(tt *args) {
				output = ""
				duc := jrpApp.NewDownloadUseCase()
				if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
					t.Errorf("Failed to download WordNet Japan DB file: %v", err)
				}
				cm := database.NewConnectionManager(proxy.NewSql())
				if err := cm.InitializeConnection(
					database.ConnectionConfig{
						DBName: database.JrpDB,
						DBType: database.SQLite,
						DSN:    filepath.Join(os.TempDir(), "jrp.db"),
					},
				); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
				if err := cm.InitializeConnection(
					database.ConnectionConfig{
						DBName: database.WNJpnDB,
						DBType: database.SQLite,
						DSN:    filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
			},
			cleanup: func() {
				rootOps.Version = false
				output = ""
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(&tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := runRoot(tt.args.cmd, tt.args.cmdArgs, tt.args.generateCmd, tt.args.versionCmd); (err != nil) != tt.wantErr {
				t.Errorf("runRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(output) == 0 {
				t.Errorf("runRoot() output is empty")
			}
			if tt.want != "" && output != tt.want {
				t.Errorf("runRoot() output = %v, want %v", output, tt.want)
			}
		})
	}
}
