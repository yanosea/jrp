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

	historyDomain "github.com/yanosea/jrp/app/domain/jrp/history"
	"github.com/yanosea/jrp/app/infrastructure/database"
	"github.com/yanosea/jrp/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewRemoveCommand(t *testing.T) {
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
			got := NewRemoveCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewRemoveCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run remove command : %v", err)
				}
			}
		})
	}
}

func Test_runRemove(t *testing.T) {
	var output string
	origRemoveOps := removeOps
	origPu := presenter.Pu

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
			name: "positive testing (not all, not no-confirm)",
			args: args{
				cmd:    &c.Command{},
				args:   []string{"1"},
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
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    color.GreenString("✅ Removed successfully!"),
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				removeOps.All = false
				removeOps.NoConfirm = false
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
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
				removeOps = origRemoveOps
				output = ""
			},
		},
		{
			name: "positive testing (all, not no-confirm, answer is y)",
			args: args{
				cmd:    nil,
				args:   []string{"1"},
				output: &output,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    color.GreenString("✅ Removed successfully!"),
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				removeOps.All = true
				removeOps.NoConfirm = false
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
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
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().Run().Return("y", nil)
				mockPromptUtil := utility.NewMockPromptUtil(mockCtrl)
				mockPromptUtil.EXPECT().GetPrompt("Proceed with removing all the histories? [y/N]").Return(mockPrompt)
				presenter.Pu = mockPromptUtil
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				removeOps = origRemoveOps
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "positive testing (all, not no-confirm, answer is not y)",
			args: args{
				cmd:    nil,
				args:   []string{"1"},
				output: &output,
			},
			testData: []*historyDomain.History{
				{
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
			want:    color.YellowString("🚫 Cancelled removing all the histories."),
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				removeOps.All = true
				removeOps.NoConfirm = false
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
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
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().Run().Return("n", nil)
				mockPromptUtil := utility.NewMockPromptUtil(mockCtrl)
				mockPromptUtil.EXPECT().GetPrompt("Proceed with removing all the histories? [y/N]").Return(mockPrompt)
				presenter.Pu = mockPromptUtil
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				removeOps = origRemoveOps
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "positive testing (no args)",
			args: args{
				cmd:    nil,
				args:   []string{},
				output: &output,
			},
			testData: nil,
			want:     color.YellowString("⚡ No ID arguments specified..."),
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *args) {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "positive testing (no histories to remove)",
			args: args{
				cmd:    nil,
				args:   []string{"1"},
				output: &output,
			},
			testData: nil,
			want:     color.YellowString("⚡ No histories to remove..."),
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *args) {
				removeOps.All = false
				removeOps.NoConfirm = false
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
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
				removeOps = origRemoveOps
				output = ""
			},
		},
		{
			name: "negative testing (strconv.Atoi(arg) failed)",
			args: args{
				cmd:    nil,
				args:   []string{"test"},
				output: &output,
			},
			testData: nil,
			want:     color.RedString("🚨 The ID argument must be an integer..."),
			wantErr:  true,
			setup: func(_ *gomock.Controller, _ *args) {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (presenter.RunPrompt() failed)",
			args: args{
				cmd:    nil,
				args:   []string{"1"},
				output: &output,
			},
			testData: nil,
			want:     "",
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				removeOps.All = true
				removeOps.NoConfirm = false
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().Run().Return("", errors.New("PromptProxy.Run() failed"))
				mockPromptUtil := utility.NewMockPromptUtil(mockCtrl)
				mockPromptUtil.EXPECT().GetPrompt("Proceed with removing all the histories? [y/N]").Return(mockPrompt)
				presenter.Pu = mockPromptUtil
				output = ""
			},
			cleanup: func() {
				removeOps = origRemoveOps
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (rhuc.Run() failed)",
			args: args{
				cmd:    nil,
				args:   []string{"1"},
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
			if err := runRemove(tt.args.cmd, tt.args.args, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *tt.args.output != tt.want {
				t.Errorf("runRemove() = %v, want %v", *tt.args.output, tt.want)
			}
		})
	}
}
