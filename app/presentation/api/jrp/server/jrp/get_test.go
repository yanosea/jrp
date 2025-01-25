package jrp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	wnjpnApp "github.com/yanosea/jrp/v2/app/application/wnjpn"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/presentation/api/jrp/formatter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestBindGetJrpHandler(t *testing.T) {
	type args struct {
		g proxy.Group
	}
	tests := []struct {
		name  string
		args  args
		setup func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name: "positive testing",
			args: args{
				g: nil,
			},
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mockGroup := proxy.NewMockGroup(mockCtrl)
				mockGroup.EXPECT().GET("/jrp", gomock.Any())
				tt.g = mockGroup
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
			BindGetJrpHandler(tt.args.g)
		})
	}
}

func Test_getJrp(t *testing.T) {
	origFormat := format
	origJu := formatter.Ju
	origFunc := database.GetConnectionManagerFunc
	origNewFetchWordsUseCase := wnjpnApp.NewFetchWordsUseCase
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type args struct {
		c echo.Context
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
				c: nil,
			},
			wantErr: false,
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
				tt.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/jrp", nil), httptest.NewRecorder())
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (connManager == nil)",
			args: args{
				c: nil,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				tt.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/jrp", nil), httptest.NewRecorder())
			},
			cleanup: nil,
		},
		{
			name: "negative testing (connectionManager.GetConnection(WNJpnDB) failed)",
			args: args{
				c: nil,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.WNJpnDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				database.GetConnectionManagerFunc = func() database.ConnectionManager {
					return mockConnManager
				}
				tt.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/jrp", nil), httptest.NewRecorder())
			},
			cleanup: func() {
				database.GetConnectionManagerFunc = origFunc
			},
		},
		{
			name: "negative testing (fwuc.Run() failed)",
			args: args{
				c: nil,
			},
			wantErr: false,
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
				mockWordQueryService := wnjpnApp.NewMockWordQueryService(mockCtrl)
				mockWordQueryService.EXPECT().
					FindByLangIsAndPosIn(gomock.Any(), "jpn", gomock.Any()).
					Return(nil, errors.New("WordQueryService.FindByLangIsAndPosIn() failed"))
				origNewFetchWordsUseCase := wnjpnApp.NewFetchWordsUseCase
				wnjpnApp.NewFetchWordsUseCase = func(wordQueryService wnjpnApp.WordQueryService) *wnjpnApp.FetchWordsUseCaseStruct {
					return origNewFetchWordsUseCase(mockWordQueryService)
				}
				tt.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/jrp", nil), httptest.NewRecorder())
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				wnjpnApp.NewFetchWordsUseCase = origNewFetchWordsUseCase
			},
		},
		{
			name: "negative testing (formatter.NewFormatter(format) failed)",
			args: args{
				c: nil,
			},
			wantErr: false,
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
				format = "test"
				tt.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/jrp", nil), httptest.NewRecorder())
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				format = origFormat
			},
		},
		{
			name: "negative testing (f.Format() failed)",
			args: args{
				c: nil,
			},
			wantErr: false,
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
				mockJu := utility.NewMockJsonUtil(mockCtrl)
				mockJu.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("JsonUtil.Marshal() failed"))
				formatter.Ju = mockJu
				tt.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/jrp", nil), httptest.NewRecorder())
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				formatter.Ju = origJu
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
			if err := getJrp(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getJrp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
