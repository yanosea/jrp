package history

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
	c "github.com/spf13/cobra"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewShowCommand(t *testing.T) {
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
			got := NewShowCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewShowCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run show command : %v", err)
				}
			}
		})
	}
}

func Test_runShow(t *testing.T) {
	var output string
	su := utility.NewStringsUtil()
	origShowOps := showOps
	origNewFormatter := formatter.NewFormatter

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
		setup    func(mockCtrl *gomock.Controller, tt *args)
		cleanup  func()
	}{
		{
			name: "positive testing",
			args: args{
				cmd:    &c.Command{},
				args:   []string{},
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
			setup: func(_ *gomock.Controller, tt *args) {
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
				args:   []string{},
				output: &output,
			},
			want:    color.YellowString("⚡ No histories found..."),
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
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
			name: "positive testing (arg is 9)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"9"},
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
				{
					ID:     2,
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
				{
					ID:     3,
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
				{
					ID:     4,
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
				{
					ID:     5,
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
				{
					ID:     6,
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
				{
					ID:     7,
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
				{
					ID:     8,
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
				{
					ID:     9,
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
				{
					ID:     10,
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
			want:    "IDPHRASEPREFIXSUFFIXISFAVORITEDCREATEDATUPDATEDAT2testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "3testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "4testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "5testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "6testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "7testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "8testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "9testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "10testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "TOTAL:9jrps!",
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
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
			name: "positive testing (arg is 10)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"10"},
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
				{
					ID:     2,
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
				{
					ID:     3,
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
				{
					ID:     4,
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
				{
					ID:     5,
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
				{
					ID:     6,
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
				{
					ID:     7,
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
				{
					ID:     8,
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
				{
					ID:     9,
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
				{
					ID:     10,
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
			want:    "IDPHRASEPREFIXSUFFIXISFAVORITEDCREATEDATUPDATEDAT1testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "2testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "3testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "4testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "5testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "6testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "7testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "8testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "9testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "10testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "TOTAL:10jrps!",
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
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
			name: "positive testing (arg is 11, number option is 12)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"11"},
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
				{
					ID:     2,
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
				{
					ID:     3,
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
				{
					ID:     4,
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
				{
					ID:     5,
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
				{
					ID:     6,
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
				{
					ID:     7,
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
				{
					ID:     8,
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
				{
					ID:     9,
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
				{
					ID:     10,
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
				{
					ID:     11,
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
				{
					ID:     12,
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
			want:    "IDPHRASEPREFIXSUFFIXISFAVORITEDCREATEDATUPDATEDAT1testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "2testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "3testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "4testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "5testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "6testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "7testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "8testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "9testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "10testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "11testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "12testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "TOTAL:12jrps!",
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				showOps.Number = 12
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
				showOps = origShowOps
				output = ""
			},
		},
		{
			name: "positive testing (arg is 12, number option is 11)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"12"},
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
				{
					ID:     2,
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
				{
					ID:     3,
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
				{
					ID:     4,
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
				{
					ID:     5,
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
				{
					ID:     6,
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
				{
					ID:     7,
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
				{
					ID:     8,
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
				{
					ID:     9,
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
				{
					ID:     10,
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
				{
					ID:     11,
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
				{
					ID:     12,
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
			want:    "IDPHRASEPREFIXSUFFIXISFAVORITEDCREATEDATUPDATEDAT1testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "2testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "3testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "4testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "5testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "6testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "7testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "8testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "9testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "10testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "11testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "12testprefixsuffix○" + now.Format("2006-01-02") + now.Format("15:04:05") + now.Format("2006-01-02") + now.Format("15:04:05") + "TOTAL:12jrps!",
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				showOps.Number = 11
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
				showOps = origShowOps
				output = ""
			},
		},
		{
			name: "negative testing (strconv.Atoi(args[0]) failed)",
			args: args{
				cmd:    nil,
				args:   []string{"test"},
				output: &output,
			},
			testData: nil,
			want:     color.RedString("🚨 The number argument must be an integer..."),
			wantErr:  true,
			setup: func(_ *gomock.Controller, _ *args) {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (ghuc.Run() failed)",
			args: args{
				cmd:    nil,
				args:   []string{},
				output: &output,
			},
			testData: nil,
			want:     "",
			wantErr:  true,
			setup: func(_ *gomock.Controller, tt *args) {
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
				args:   []string{},
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
			setup: func(_ *gomock.Controller, tt *args) {
				showOps.Format = "test"
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
				showOps = origShowOps
				output = ""
			},
		},
		{
			name: "negative testing (f.Format() failed)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{},
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
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
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
				mockFormatter := formatter.NewMockFormatter(mockCtrl)
				mockFormatter.EXPECT().Format(gomock.Any()).Return("", errors.New("format error"))
				formatter.NewFormatter = func(format string) (formatter.Formatter, error) {
					return mockFormatter, nil
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
				formatter.NewFormatter = origNewFormatter
				output = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.args)
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
			if err := runShow(tt.args.cmd, tt.args.args, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runShow() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.testData != nil && !tt.wantErr {
				output = su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(output)))
			}
			if output != tt.want {
				t.Errorf("runShow() = %v, want %v", output, tt.want)
			}
		})
	}
}
