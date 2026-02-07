package generate

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	wnjpnApp "github.com/yanosea/jrp/v2/app/application/wnjpn"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewGenerateCommand(t *testing.T) {
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type args struct {
		cobra          proxy.Cobra
		interactiveCmd proxy.Command
		output         *string
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
				cobra: proxy.NewCobra(),
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(),
					new(string),
				),
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
				if err := cm.InitializeConnection(
					database.ConnectionConfig{
						DBName: database.WNJpnDB,
						DBType: database.SQLite,
						DSN:    filepath.Join(os.TempDir(), "wnjpn.db"),
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
			got := NewGenerateCommand(tt.args.cobra, tt.args.interactiveCmd, tt.args.output)
			if got == nil {
				t.Errorf("NewGenerateCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run the generate command : %v", err)
				}
			}
		})
	}
}

func Test_runGenerate(t *testing.T) {
	var output string
	origGenerateOps := GenerateOps
	origKu := presenter.Ku
	origFunc := database.GetConnectionManagerFunc
	origNewFetchWordsUseCase := wnjpnApp.NewFetchWordsUseCase
	origNewFormatter := formatter.NewFormatter
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type args struct {
		cmd            *c.Command
		args           []string
		interactiveCmd proxy.Command
		output         *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
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
			name: "positive testing (arg is 2)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{"2"},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
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
			name: "positive testing (arg is 2, number option is 3)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{"2"},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				GenerateOps.Number = 3
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				GenerateOps = origGenerateOps
				output = ""
			},
		},
		{
			name: "positive testing (arg is 3, number option is 2)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{"3"},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				GenerateOps.Number = 2
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				GenerateOps = origGenerateOps
				output = ""
			},
		},
		{
			name: "positive testing (prefix option is set)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				GenerateOps.Prefix = "prefix"
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				GenerateOps = origGenerateOps
				output = ""
			},
		},
		{
			name: "positive testing (suffix option is set)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				GenerateOps.Suffix = "suffix"
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				GenerateOps = origGenerateOps
				output = ""
			},
		},
		{
			name: "positive testing (interactive option is set)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				GenerateOps.Interactive = true
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return(",", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
				presenter.Ku = mockKeyboardUtil
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
				GenerateOps = origGenerateOps
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (connManager == nil)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (connManager.GetConnection(WNJpnDB) == connection not initialized)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				output = ""
			},
		},
		{
			name: "negative testing (connectionManager.GetConnection(WNJpnDB) failed)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.WNJpnDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				database.GetConnectionManagerFunc = func() database.ConnectionManager {
					return mockConnManager
				}
				output = ""
			},
			cleanup: func() {
				database.GetConnectionManagerFunc = origFunc
				output = ""
			},
		},
		{
			name: "negative testing (both prefix and suffix options are set)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				GenerateOps.Prefix = "prefix"
				GenerateOps.Suffix = "suffix"
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				GenerateOps = origGenerateOps
				output = ""
			},
		},
		{
			name: "negative testing (fwuc.Run() failed)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
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
				if err := cm.InitializeConnection(
					database.ConnectionConfig{
						DBName: database.WNJpnDB,
						DBType: database.SQLite,
						DSN:    filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
				mockWordQueryService := wnjpnApp.NewMockWordQueryService(mockCtrl)
				mockWordQueryService.EXPECT().
					FindByLangIsAndPosIn(gomock.Any(), "jpn", gomock.Any()).
					Return(nil, errors.New("WordQueryService.FindByLangIsAndPosIn() failed"))
				origNewFetchWordsUseCase := wnjpnApp.NewFetchWordsUseCase
				wnjpnApp.NewFetchWordsUseCase = func(wordQueryService wnjpnApp.WordQueryService) *wnjpnApp.FetchWordsUseCaseStruct {
					return origNewFetchWordsUseCase(mockWordQueryService)
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
				wnjpnApp.NewFetchWordsUseCase = origNewFetchWordsUseCase
				output = ""
			},
		},
		{
			name: "negative testing (strconv.Atoi(arg[0]) failed)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{"test"},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: true,
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
			name: "negative testing (shuc.Run() failed)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: true,
			setup: func(_ *gomock.Controller, tt *args) {
				GenerateOps.Suffix = "suffix"
				cm := database.NewConnectionManager(proxy.NewSql())
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				GenerateOps = origGenerateOps
				output = ""
			},
		},
		{
			name: "negative testing (formatter.NewFormatter(GenerateOps.Format) failed)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				GenerateOps.Format = "test"
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
				output = ""
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				GenerateOps = origGenerateOps
				output = ""
			},
		},
		{
			name: "negative testing (f.Format() failed)",
			args: args{
				cmd:            &c.Command{},
				args:           []string{},
				interactiveCmd: NewInteractiveCommand(proxy.NewCobra(), &output),
				output:         &output,
			},
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
				if err := cm.InitializeConnection(
					database.ConnectionConfig{
						DBName: database.WNJpnDB,
						DBType: database.SQLite,
						DSN:    filepath.Join(os.TempDir(), "wnjpn.db"),
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
			err := runGenerate(tt.args.cmd, tt.args.args, tt.args.interactiveCmd, tt.args.output)
			if tt.wantErr {
				if err == nil {
					t.Errorf("runGenerate() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("runGenerate() error = %v", err)
				}
			}
			if output != "" {
				t.Logf("runGenerate() = %s", "\n"+output)
			}
			if err != nil {
				t.Logf("runGenerate() error = %v", err)
			}
		})
	}
}
