package config

import (
	"errors"
	"reflect"
	"testing"

	baseConfig "github.com/yanosea/jrp/v2/app/config"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewJrpServerConfigurator(t *testing.T) {
	envconfig := proxy.NewEnvconfig()
	fileUtil := utility.NewFileUtil(
		proxy.NewGzip(),
		proxy.NewIo(),
		proxy.NewOs(),
	)

	type args struct {
		envconfigProxy proxy.Envconfig
		fileUtil       utility.FileUtil
	}
	tests := []struct {
		name string
		args args
		want JrpServerConfigurator
	}{
		{
			name: "positive testing",
			args: args{
				envconfigProxy: envconfig,
				fileUtil:       fileUtil,
			},
			want: &ServerConfigurator{
				BaseConfigurator: baseConfig.NewConfigurator(
					envconfig,
					fileUtil,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJrpServerConfigurator(tt.args.envconfigProxy, tt.args.fileUtil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJrpServerConfigurator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerConfigurator_GetConfig(t *testing.T) {
	type fields struct {
		BaseConfigurator *baseConfig.BaseConfigurator
	}
	tests := []struct {
		name    string
		fields  fields
		want    *JrpServerConfig
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				BaseConfigurator: &baseConfig.BaseConfigurator{
					Envconfig: nil,
					FileUtil:  nil,
				}},
			want: &JrpServerConfig{
				JrpConfig: baseConfig.JrpConfig{
					WNJpnDBType: database.SQLite,
					WNJpnDBDsn:  "~/.local/share/jrp/wnjpn.db",
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
				mockEnvconfig.EXPECT().Process("", gomock.Any()).DoAndReturn(
					func(_ string, cfg *envConfig) error {
						cfg.WnJpnDBType = database.SQLite
						cfg.WnJpnDBDsn = "XDG_DATA_HOME/jrp/wnjpn.db"
						return nil
					})
				mockFileUtil := utility.NewMockFileUtil(mockCtrl)
				mockFileUtil.EXPECT().GetXDGDataHome().Return("~/.local/share", nil)
				mockFileUtil.EXPECT().MkdirIfNotExist("~/.local/share/jrp").Return(nil)
				tt.BaseConfigurator.Envconfig = mockEnvconfig
				tt.BaseConfigurator.FileUtil = mockFileUtil
			},
		},
		{
			name: "negative testing (c.Envconfig.Process(\"\", &config) failed)",
			fields: fields{
				BaseConfigurator: &baseConfig.BaseConfigurator{
					Envconfig: nil,
					FileUtil:  nil,
				}},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
				mockEnvconfig.EXPECT().Process("", gomock.Any()).Return(errors.New("Envconfig.Process() failed"))
				tt.BaseConfigurator.Envconfig = mockEnvconfig
			},
		},
		{
			name: "negative testing (c.FileUtil.GetXDGDataHome() failed)",
			fields: fields{
				BaseConfigurator: &baseConfig.BaseConfigurator{
					Envconfig: nil,
					FileUtil:  nil,
				}},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
				mockEnvconfig.EXPECT().Process("", gomock.Any()).DoAndReturn(
					func(_ string, cfg *envConfig) error {
						cfg.WnJpnDBType = database.SQLite
						cfg.WnJpnDBDsn = "XDG_DATA_HOME/jrp/wnjpn.db"
						return nil
					})
				mockFileUtil := utility.NewMockFileUtil(mockCtrl)
				mockFileUtil.EXPECT().GetXDGDataHome().Return("", errors.New("FileUtil.GetXDGDataHome() failed"))
				tt.BaseConfigurator.Envconfig = mockEnvconfig
				tt.BaseConfigurator.FileUtil = mockFileUtil
			},
		},
		{
			name: "negative testing (c.FileUtil.MkdirIfNotExist(filepath.Dir(config.WNJpnDBConnectionString)) failed)",
			fields: fields{
				BaseConfigurator: &baseConfig.BaseConfigurator{
					Envconfig: nil,
					FileUtil:  nil,
				},
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
				mockEnvconfig.EXPECT().Process("", gomock.Any()).DoAndReturn(
					func(_ string, cfg *envConfig) error {
						cfg.WnJpnDBType = database.SQLite
						cfg.WnJpnDBDsn = "XDG_DATA_HOME/jrp/wnjpn.db"
						return nil
					})
				mockFileUtil := utility.NewMockFileUtil(mockCtrl)
				mockFileUtil.EXPECT().GetXDGDataHome().Return("~/.local/share", nil)
				mockFileUtil.EXPECT().MkdirIfNotExist("~/.local/share/jrp").Return(errors.New("FileUtil.MkdirIfNotExist() failed"))
				tt.BaseConfigurator.Envconfig = mockEnvconfig
				tt.BaseConfigurator.FileUtil = mockFileUtil
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
			c := &ServerConfigurator{
				BaseConfigurator: tt.fields.BaseConfigurator,
			}
			got, err := c.GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerConfigurator.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerConfigurator.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
