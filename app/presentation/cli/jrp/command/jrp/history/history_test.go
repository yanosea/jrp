package history

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	c "github.com/spf13/cobra"

	historyDomain "github.com/yanosea/jrp/app/domain/jrp/history"
	"github.com/yanosea/jrp/app/infrastructure/database"
	"github.com/yanosea/jrp/app/infrastructure/jrp/repository"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"
)

var (
	now = time.Now()
)

func TestNewHistoryCommand(t *testing.T) {
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
			got := NewHistoryCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewHistoryCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run history command : %v", err)
				}
			}
		})
	}
}

func Test_runHistory(t *testing.T) {
	var output string
	su := utility.NewStringsUtil()

	type args struct {
		cmd     *c.Command
		showCmd proxy.Command
		args    []string
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
				cmd: nil,
				showCmd: NewShowCommand(
					proxy.NewCobra(),
					&output,
				),
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
			if err := runHistory(tt.args.cmd, tt.args.showCmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(output))) != tt.want {
				t.Errorf("runHistory() = %v, want %v", su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(output))), tt.want)
			}
		})
	}
}
