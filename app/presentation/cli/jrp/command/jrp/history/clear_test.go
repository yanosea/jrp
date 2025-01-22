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

func TestNewClearCommand(t *testing.T) {
	origClearOps := clearOps
	origPu := presenter.Pu

	type args struct {
		cobra  proxy.Cobra
		output *string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(mockCtrl *gomock.Controller, tt *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				cobra:  proxy.NewCobra(),
				output: new(string),
			},
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				clearOps.Force = false
				clearOps.NoConfirm = true
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().Run().Return("n", nil)
				mockPromptUtil := utility.NewMockPromptUtil(mockCtrl)
				mockPromptUtil.EXPECT().GetPrompt("Proceed with clearing the histories? [y/N]").Return(mockPrompt)
				presenter.Pu = mockPromptUtil
			},
			cleanup: func() {
				clearOps = origClearOps
				presenter.Pu = origPu
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
			got := NewClearCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewClearCommand() = %v, want not nil", got)
			} else {
				if err := got.RunE(nil, []string{}); err != nil {
					t.Errorf("Failed to run the clear command: %v", err)
				}
			}
		})
	}
}

func Test_runClear(t *testing.T) {
	var output string
	origClearOps := clearOps
	origPu := presenter.Pu

	type args struct {
		cmd    *c.Command
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
			name: "positive testing (not force, not no-confirm)",
			args: args{
				cmd:    nil,
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
			want:    color.GreenString("✅ Cleared successfully!"),
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				clearOps.Force = false
				clearOps.NoConfirm = false
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
				mockPromptUtil.EXPECT().GetPrompt("Proceed with clearing the histories? [y/N]").Return(mockPrompt)
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
				output = ""
			},
		},
		{
			name: "positive testing (not force, no-confirm)",
			args: args{
				cmd:    nil,
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
			want:    color.GreenString("✅ Cleared successfully!"),
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				clearOps.Force = false
				clearOps.NoConfirm = true
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
				clearOps.NoConfirm = true
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				clearOps = origClearOps
				output = ""
			},
		},
		{
			name: "positive testing (force, not no-confirm, answer is y)",
			args: args{
				cmd:    nil,
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
			want:    color.GreenString("✅ Cleared successfully!"),
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				clearOps.Force = true
				clearOps.NoConfirm = false
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
				mockPromptUtil.EXPECT().GetPrompt("Proceed with clearing the histories? [y/N]").Return(mockPrompt)
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
				clearOps = origClearOps
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "positive testing (force, not no-confirm, answer is not y)",
			args: args{
				cmd:    nil,
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
			want:    color.YellowString("🚫 Cancelled clearing the histories."),
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				clearOps.Force = true
				clearOps.NoConfirm = false
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
				mockPromptUtil.EXPECT().GetPrompt("Proceed with clearing the histories? [y/N]").Return(mockPrompt)
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
				clearOps = origClearOps
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "positive testing (no histories to clear)",
			args: args{
				cmd:    nil,
				output: &output,
			},
			testData: nil,
			want:     color.YellowString("⚡ No histories to clear..."),
			wantErr:  false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				clearOps.Force = false
				clearOps.NoConfirm = false
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
				mockPromptUtil.EXPECT().GetPrompt("Proceed with clearing the histories? [y/N]").Return(mockPrompt)
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
				output = ""
			},
		},
		{
			name: "negative testing (presenter.RunPrompt() failed)",
			args: args{
				cmd:    nil,
				output: &output,
			},
			testData: nil,
			want:     "",
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				clearOps.Force = true
				clearOps.NoConfirm = false
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().Run().Return("", errors.New("PromptProxy.Run() failed"))
				mockPromptUtil := utility.NewMockPromptUtil(mockCtrl)
				mockPromptUtil.EXPECT().GetPrompt("Proceed with clearing the histories? [y/N]").Return(mockPrompt)
				presenter.Pu = mockPromptUtil
				output = ""
			},
			cleanup: func() {
				clearOps = origClearOps
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (rhuc.Run() failed)",
			args: args{
				cmd:    nil,
				output: &output,
			},
			testData: nil,
			want:     "",
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				clearOps.Force = true
				clearOps.NoConfirm = false
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				database.NewConnectionManager(proxy.NewSql())
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().Run().Return("y", nil)
				mockPromptUtil := utility.NewMockPromptUtil(mockCtrl)
				mockPromptUtil.EXPECT().GetPrompt("Proceed with clearing the histories? [y/N]").Return(mockPrompt)
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
				output = ""
				presenter.Pu = origPu
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
			if err := runClear(tt.args.cmd, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runClear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *tt.args.output != tt.want {
				t.Errorf("runClear() = %v, want %v", *tt.args.output, tt.want)
			}
		})
	}
}
