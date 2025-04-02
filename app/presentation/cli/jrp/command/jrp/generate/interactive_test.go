package generate

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func TestNewInteractiveCommand(t *testing.T) {
	origKu := presenter.Ku
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type args struct {
		cobra  proxy.Cobra
		output *string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				cobra:  proxy.NewCobra(),
				output: new(string),
			},
			setup: func(mockCtrl *gomock.Controller) {
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
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				presenter.Ku = origKu
			},
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
			got := NewInteractiveCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewInteractiveCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run the interactive command : %v", err)
				}
			}
		})
	}
}

func Test_runInteractive(t *testing.T) {
	var output string
	origInteractiveOps := interactiveOps
	origKu := presenter.Ku
	origFunc := database.GetConnectionManagerFunc
	origNewFetchWordsUseCase := wnjpnApp.NewFetchWordsUseCase
	origPrint := presenter.Print
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type args struct {
		cmd    *c.Command
		output *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
		cleanup func()
	}{
		{
			name: "positive testing (keyboard input: \"u\")",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("u", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "positive testing (keyboard input: \"i\")",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("i", nil)
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "positive testing (keyboard input: \"j\")",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("j", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "positive testing (keyboard input: \"k\")",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("k", nil)
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "positive testing (keyboard input: \"m\")",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("m", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "positive testing (keyboard input: \",\")",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "positive testing (prefix option is set)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				interactiveOps.Prefix = "prefix"
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
				interactiveOps = origInteractiveOps
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "positive testing (suffix option is set)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				interactiveOps.Suffix = "suffix"
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
				interactiveOps = origInteractiveOps
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (connManager == nil)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				cmd:    &c.Command{},
				output: &output,
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
				cmd:    &c.Command{},
				output: &output,
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
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				interactiveOps.Prefix = "prefix"
				interactiveOps.Suffix = "suffix"
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
				interactiveOps = origInteractiveOps
				output = ""
			},
		},
		{
			name: "negative testing (fwuc.Run() failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
			name: "negative testing (presenter.Print(os.Stdout, formatter.Blue(\"🔄 Phase : \"+strconv.Itoa(phase)))) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Phase") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (formatter.NewFormatter(interactiveOps.Format) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				interactiveOps.Format = "test"
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
				interactiveOps = origInteractiveOps
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\")) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				printCount := 0
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Phase") {
						return origPrint(writer, output)
					}
					if output == "\n" {
						printCount++
						if printCount == 1 {
							return errors.New("Print(\\n) failed")
						}
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, o)) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				printCount := 0
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Phase") || output == "\n" {
						return origPrint(writer, output)
					}
					printCount++
					if printCount == 1 {
						return errors.New("Print(o) failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\"))) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				newlineCount := 0
				presenter.Print = func(writer io.Writer, output string) error {
					if output == "\n" {
						newlineCount++
						if newlineCount == 2 {
							return errors.New("Print(\\n) failed")
						}
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, formatter.Yellow(interactivePromptLabel))) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Press either key below") {
						return errors.New("Print(interactivePromptLabel) failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.OpenKeyboard() failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(errors.New("KeyboardUtil.OpenKeyboard() failed"))
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
				interactiveOps = origInteractiveOps
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.GetKey(interactiveOps.Timeout) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("", errors.New("KeyboardUtil.GetKey() failed"))
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
				interactiveOps = origInteractiveOps
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.CloseKeyboard() failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil).Times(1)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("m", nil).Times(1)
				mockKeyboardUtil.EXPECT().CloseKeyboard().Return(errors.New("KeyboardUtil.CloseKeyboard() failed")).Times(1)
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (shuc.Run() failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("u", nil)
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
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, formatter.Green(\"✅ Favorited successfully!\"))) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("u", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
				presenter.Ku = mockKeyboardUtil
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Favorited successfully") {
						return errors.New("Print(Favorited successfully) failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\"))) after favorited failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("u", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
				presenter.Ku = mockKeyboardUtil
				var successMsgShown bool
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Favorited successfully") {
						successMsgShown = true
						return origPrint(writer, output)
					}
					if output == "\n" && successMsgShown {
						return errors.New("Print(\\n) after favorited failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, formatter.Green(\"✅ Saved successfully!\")) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("j", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
				presenter.Ku = mockKeyboardUtil
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Saved successfully") {
						return errors.New("Print(Saved successfully) failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\"))) after saved failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return("j", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
				presenter.Ku = mockKeyboardUtil
				var successMsgShown bool
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Saved successfully") {
						successMsgShown = true
						return origPrint(writer, output)
					}
					if output == "\n" && successMsgShown {
						return errors.New("Print(\\n) after saved failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, formatter.Yellow(\"⏩ Skip!\")) failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				var printCalls int
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil).Times(1)
				mockKeyboardUtil.EXPECT().GetKey(gomock.Any()).Return("m", nil).Times(1)
				mockKeyboardUtil.EXPECT().CloseKeyboard().Times(1)
				presenter.Ku = mockKeyboardUtil
				presenter.Print = func(writer io.Writer, output string) error {
					printCalls++
					if output == formatter.Yellow("⏩ Skip!") {
						return errors.New("Print(Skip) failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\")) after Skip failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil).Times(1)
				mockKeyboardUtil.EXPECT().GetKey(gomock.Any()).Return("m", nil).Times(1)
				mockKeyboardUtil.EXPECT().CloseKeyboard().Times(1)
				presenter.Ku = mockKeyboardUtil
				var skipMsgShown bool
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Skip") {
						skipMsgShown = true
						return origPrint(writer, output)
					}
					if output == "\n" && skipMsgShown {
						return errors.New("Print(\\n) after skip failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				presenter.Ku = origKu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"🚪 Exit!\") failed)",
			args: args{
				cmd:    &c.Command{},
				output: &output,
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
				mockKeyboardUtil := utility.NewMockKeyboardUtil(mockCtrl)
				mockKeyboardUtil.EXPECT().OpenKeyboard().Return(nil)
				mockKeyboardUtil.EXPECT().GetKey(interactiveOps.Timeout).Return(",", nil)
				mockKeyboardUtil.EXPECT().CloseKeyboard()
				presenter.Ku = mockKeyboardUtil
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "Exit") {
						return errors.New("Print(Exit) failed")
					}
					return origPrint(writer, output)
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
				presenter.Print = origPrint
				presenter.Ku = origKu
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
			err := runInteractive(tt.args.cmd, tt.args.output)
			if tt.wantErr {
				if err == nil {
					t.Errorf("runInteractive() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("runInteractive() error = %v", err)
				}
			}
			if output != "" {
				t.Logf("runInteractive() = %s", "\n"+output)
			}
			if err != nil {
				t.Logf("runInteractive() error = %v", err)
			}
		})
	}
}
