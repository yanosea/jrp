package history

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
	c "github.com/spf13/cobra"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/infrastructure/jrp/repository"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

func TestNewSearchCommand(t *testing.T) {
	type args struct {
		cobra  proxy.Cobra
		output *string
	}
	tests := []struct {
		name    string
		args    args
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				cobra:  proxy.NewCobra(),
				output: new(string),
			},
			setup: func() {
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
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			got := NewSearchCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewSearchCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run search command : %v", err)
				}
			}
		})
	}
}

func Test_runSearch(t *testing.T) {
	var output string
	su := utility.NewStringsUtil()
	origSearchOps := searchOps

	type args struct {
		cmd    *c.Command
		args   []string
		output *string
	}
	tests := []struct {
		name     string
		args     args
		testData []*historyDomain.History
		want     string
		wantErr  bool
		setup    func(tt *args)
		cleanup  func()
	}{
		{
			name: "positive testing",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"test"},
				output: &output,
			},
			testData: []*historyDomain.History{
				{
					ID:     1,
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    "IDPHRASEPREFIXSUFFIXISFAVORITEDCREATEDATUPDATEDAT1testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "TOTAL:1jrps!",
			wantErr: false,
			setup: func(tt *args) {
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
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				output = ""
			},
		},
		{
			name: "positive testing (no histories in the database)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"test"},
				output: &output,
			},
			want:    color.YellowString("⚡ No histories found..."),
			wantErr: false,
			setup: func(tt *args) {
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
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				output = ""
			},
		},
		{
			name: "positive testing (args are empty)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{},
				output: &output,
			},
			testData: nil,
			want:     color.YellowString("⚡ No keywords provided..."),
			wantErr:  false,
			setup: func(tt *args) {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (shuc.Run() failed)",
			args: args{
				cmd:    nil,
				args:   []string{"test"},
				output: &output,
			},
			testData: nil,
			want:     "",
			wantErr:  true,
			setup: func(tt *args) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				database.NewConnectionManager(proxy.NewSql())
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				output = ""
			},
		},
		{
			name: "negative testing (formatter.NewFormatter(showOps.Format) failed)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"test"},
				output: &output,
			},
			testData: []*historyDomain.History{
				{
					ID:     1,
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    color.RedString("❌ Failed to create a formatter..."),
			wantErr: true,
			setup: func(tt *args) {
				searchOps.Format = "test"
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
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				searchOps = origSearchOps
				output = ""
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
			h := repository.NewHistoryRepository()
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(context.Background(), tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			if err := runSearch(tt.args.cmd, tt.args.args, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.testData != nil && !tt.wantErr {
				output = su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(output)))
			}
			if output != tt.want {
				t.Errorf("runSearch() = %v, want %v", output, tt.want)
			}
		})
	}
}
