package server

import (
	"errors"
	o "os"
	"path/filepath"
	"reflect"
	"testing"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func Test_newServer(t *testing.T) {
	echos := proxy.NewEchos()

	type args struct {
		echos proxy.Echos
	}
	tests := []struct {
		name string
		args args
		want Server
	}{
		{
			name: "positeive testing",
			args: args{
				echos: echos,
			},
			want: &server{
				ConnectionManager: nil,
				Echos:             echos,
				Logger:            nil,
				Port:              "",
				Route:             nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newServer(tt.args.echos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_Init(t *testing.T) {
	echos := proxy.NewEchos()
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(o.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type fields struct {
		ConnectionManager database.ConnectionManager
		Echos             proxy.Echos
		Logger            proxy.Logger
		Port              string
		Route             proxy.Echo
	}
	type args struct {
		envconfig proxy.Envconfig
		fileUtil  utility.FileUtil
		sql       proxy.Sql
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		setup   func(mockCtrl *gomock.Controller, ta *args, tf *fields)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				ConnectionManager: nil,
				Echos:             echos,
				Logger:            nil,
				Port:              "",
				Route:             nil,
			},
			args: args{
				envconfig: proxy.NewEnvconfig(),
				fileUtil: utility.NewFileUtil(
					proxy.NewGzip(),
					proxy.NewIo(),
					proxy.NewOs(),
				),
				sql: proxy.NewSql(),
			},
			want: 0,
			setup: func(_ *gomock.Controller, _ *args, _ *fields) {
				if err := o.Setenv("JRP_SERVER_PORT", "8080"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_PORT"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
		{
			name: "negative testing (configurator.GetConfig() failed)",
			fields: fields{
				ConnectionManager: nil,
				Echos:             nil,
				Logger:            nil,
				Port:              "",
				Route:             nil,
			},
			args: args{
				envconfig: nil,
				fileUtil: utility.NewFileUtil(
					proxy.NewGzip(),
					proxy.NewIo(),
					proxy.NewOs(),
				),
				sql: proxy.NewSql(),
			},
			want: 1,
			setup: func(mockCtrl *gomock.Controller, ta *args, tf *fields) {
				if err := o.Setenv("JRP_SERVER_PORT", "8080"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				mockGroup := proxy.NewMockGroup(mockCtrl)
				mockGroup.EXPECT().GET(gomock.Any(), gomock.Any())
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Group(gomock.Any()).Return(mockGroup)
				mockEcho.EXPECT().Get("/swagger/*", gomock.Any())
				mockLogger := proxy.NewMockLogger(mockCtrl)
				mockLogger.EXPECT().Fatal(gomock.Any())
				mockEchos := proxy.NewMockEchos(mockCtrl)
				mockEchos.EXPECT().NewEcho().Return(mockEcho, mockLogger)
				tf.Echos = mockEchos
				mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
				mockEnvconfig.EXPECT().Process("", gomock.Any()).Return(errors.New("EnvconfigProxy.Process() failed"))
				ta.envconfig = mockEnvconfig
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_PORT"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
		{
			name: "negative testing (conf.WNJpnDBType == database.SQLite && !fileUtil.IsExist(conf.WNJpnDBDsn))",
			fields: fields{
				ConnectionManager: nil,
				Echos:             echos,
				Logger:            nil,
				Port:              "",
				Route:             nil,
			},
			args: args{
				envconfig: proxy.NewEnvconfig(),
				fileUtil:  nil,
				sql:       proxy.NewSql(),
			},
			want: 1,
			setup: func(mockCtrl *gomock.Controller, ta *args, tf *fields) {
				if err := o.Setenv("JRP_SERVER_PORT", "8080"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				mockGroup := proxy.NewMockGroup(mockCtrl)
				mockGroup.EXPECT().GET(gomock.Any(), gomock.Any())
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Group(gomock.Any()).Return(mockGroup)
				mockEcho.EXPECT().Get("/swagger/*", gomock.Any())
				mockLogger := proxy.NewMockLogger(mockCtrl)
				mockLogger.EXPECT().Fatal(gomock.Any())
				mockEchos := proxy.NewMockEchos(mockCtrl)
				mockEchos.EXPECT().NewEcho().Return(mockEcho, mockLogger)
				tf.Echos = mockEchos
				mockFileUtil := utility.NewMockFileUtil(mockCtrl)
				mockFileUtil.EXPECT().GetXDGDataHome().Return("~/.local/share", nil)
				mockFileUtil.EXPECT().MkdirIfNotExist(gomock.Any()).Return(nil)
				mockFileUtil.EXPECT().IsExist(filepath.Join(o.TempDir(), "wnjpn.db")).Return(false)
				ta.fileUtil = mockFileUtil
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_PORT"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
		{
			name: "negative testing (WNJpnDB InitializeConnection failed)",
			fields: fields{
				ConnectionManager: nil,
				Echos:             echos,
				Logger:            nil,
				Port:              "",
				Route:             nil,
			},
			args: args{
				envconfig: proxy.NewEnvconfig(),
				fileUtil: utility.NewFileUtil(
					proxy.NewGzip(),
					proxy.NewIo(),
					proxy.NewOs(),
				),
				sql: proxy.NewSql(),
			},
			want: 1,
			setup: func(mockCtrl *gomock.Controller, ta *args, tf *fields) {
				if err := o.Setenv("JRP_SERVER_PORT", "8080"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB_TYPE", "sqlite"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("JRP_SERVER_WNJPN_DB", filepath.Join(o.TempDir(), "wnjpn.db")); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				mockGroup := proxy.NewMockGroup(mockCtrl)
				mockGroup.EXPECT().GET(gomock.Any(), gomock.Any())
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Group(gomock.Any()).Return(mockGroup)
				mockEcho.EXPECT().Get("/swagger/*", gomock.Any())
				mockLogger := proxy.NewMockLogger(mockCtrl)
				mockLogger.EXPECT().Fatal(gomock.Any())
				mockEchos := proxy.NewMockEchos(mockCtrl)
				mockEchos.EXPECT().NewEcho().Return(mockEcho, mockLogger)
				tf.Echos = mockEchos
				mockConnectionManager := database.NewMockConnectionManager(mockCtrl)
				mockConnectionManager.EXPECT().InitializeConnection(gomock.Any()).Return(errors.New("ConnectionManager.InitializeConnection() failed"))
				tf.ConnectionManager = mockConnectionManager
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_PORT"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB_TYPE"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("JRP_SERVER_WNJPN_DB"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.args, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			s := &server{
				ConnectionManager: tt.fields.ConnectionManager,
				Echos:             tt.fields.Echos,
				Logger:            tt.fields.Logger,
				Port:              tt.fields.Port,
				Route:             tt.fields.Route,
			}
			if got := s.Init(tt.args.envconfig, tt.args.fileUtil, tt.args.sql); got != tt.want {
				t.Errorf("server.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_Run(t *testing.T) {
	type fields struct {
		ConnectionManager database.ConnectionManager
		Echos             proxy.Echos
		Logger            proxy.Logger
		Port              string
		Route             proxy.Echo
	}
	tests := []struct {
		name         string
		fields       fields
		wantExitCode int
		setup        func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				ConnectionManager: nil,
				Echos:             proxy.NewEchos(),
				Logger:            nil,
				Port:              "8080",
				Route:             nil,
			},
			wantExitCode: 0,
			setup: func(mockCtrl *gomock.Controller, tf *fields) {
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Start(":" + tf.Port).Return(nil)
				tf.Route = mockEcho
			},
		},
		{
			name: "negative testing (s.Route.Start(\":\" + s.Port) failed)",
			fields: fields{
				ConnectionManager: nil,
				Echos:             proxy.NewEchos(),
				Logger:            nil,
				Port:              "8080",
				Route:             nil,
			},
			wantExitCode: 1,
			setup: func(mockCtrl *gomock.Controller, tf *fields) {
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Start(":" + tf.Port).Return(errors.New("EchoProxy.Start() failed"))
				mockLogger := proxy.NewMockLogger(mockCtrl)
				mockLogger.EXPECT().Fatal(gomock.Any())
				tf.Route = mockEcho
				tf.Logger = mockLogger
			},
		},
		{
			name: "negative testing (s.ConnectionManager.CloseAllConnections() failed)",
			fields: fields{
				ConnectionManager: nil,
				Echos:             proxy.NewEchos(),
				Logger:            nil,
				Port:              "8080",
				Route:             nil,
			},
			wantExitCode: 1,
			setup: func(mockCtrl *gomock.Controller, tf *fields) {
				mockConnectionManager := database.NewMockConnectionManager(mockCtrl)
				mockConnectionManager.EXPECT().CloseAllConnections().Return(errors.New("ConnectionManager.CloseAllConnections() failed"))
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Start(":" + tf.Port).Return(errors.New("EchoProxy.Start() failed"))
				mockLogger := proxy.NewMockLogger(mockCtrl)
				mockLogger.EXPECT().Fatal(gomock.Any())
				mockLogger.EXPECT().Fatal(gomock.Any())
				tf.ConnectionManager = mockConnectionManager
				tf.Route = mockEcho
				tf.Logger = mockLogger
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			s := &server{
				ConnectionManager: tt.fields.ConnectionManager,
				Echos:             tt.fields.Echos,
				Logger:            tt.fields.Logger,
				Port:              tt.fields.Port,
				Route:             tt.fields.Route,
			}
			if gotExitCode := s.Run(); gotExitCode != tt.wantExitCode {
				t.Errorf("server.Run() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
		})
	}
}
