package testutility

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/user"

	"github.com/yanosea/jrp/mock/app/proxy/os"
	"go.uber.org/mock/gomock"
)

func TestNewTestEnvSetter(t *testing.T) {
	osProxy := osproxy.New()

	type args struct {
		osProxy osproxy.Os
	}
	tests := []struct {
		name string
		args args
		want *TestEnvSetter
	}{
		{
			name: "positive testing",
			args: args{
				osProxy: osProxy,
			},
			want: &TestEnvSetter{
				OsProxy: osProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTestEnvSetter(tt.args.osProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTestEnvSetter() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestTestEnvSetter_SetTestEnv(t *testing.T) {
	osProxy := osproxy.New()
	dbFilePathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	defaultWnJpnDBFileDirPath, err := dbFilePathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFilePathProvider.GetWNJpnDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
	}
	defaultJrpDBFileDirPath, err := dbFilePathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFilePathProvider.GetJrpDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
	}
	testEnvSetter := NewTestEnvSetter(osProxy)

	type fields struct {
		OsProxy osproxy.Os
	}
	tests := []struct {
		name                   string
		fields                 fields
		wantWnJpnDBFileDirPath string
		wantJrpDBFileDirPath   string
		wantErr                bool
		setup                  func(mockCtrl *gomock.Controller, tt *fields)
		cleanup                func()
	}{
		{
			name: "positive testing",
			fields: fields{
				OsProxy: osproxy.New(),
			},
			wantWnJpnDBFileDirPath: osProxy.TempDir(),
			wantJrpDBFileDirPath:   osProxy.TempDir(),
			wantErr:                false,
			setup:                  nil,
			cleanup: func() {
				if err := testEnvSetter.UnsetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.UnsetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
		},
		{
			name: "negative testing (OsProxy.Setenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR) failed)",
			fields: fields{
				OsProxy: nil,
			},
			wantWnJpnDBFileDirPath: defaultWnJpnDBFileDirPath,
			wantJrpDBFileDirPath:   defaultJrpDBFileDirPath,
			wantErr:                true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().TempDir().Return(osProxy.TempDir())
				mockOsProxy.EXPECT().Setenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR, osProxy.TempDir()).Return(errors.New("OsProxy.Setenv() failed"))
				fields.OsProxy = mockOsProxy
			},
			cleanup: func() {
				if err := testEnvSetter.UnsetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.UnsetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
		},
		{
			name: "negative testing (OsProxy.Setenv(dbfiledirpathprovider.JRP_ENV_JRP_DB_FILE_DIR) failed)",
			fields: fields{
				OsProxy: nil,
			},
			wantWnJpnDBFileDirPath: defaultWnJpnDBFileDirPath,
			wantJrpDBFileDirPath:   defaultJrpDBFileDirPath,
			wantErr:                true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().TempDir().Return(osProxy.TempDir())
				mockOsProxy.EXPECT().Setenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR, osProxy.TempDir()).Return(nil)
				mockOsProxy.EXPECT().Setenv(dbfiledirpathprovider.JRP_ENV_JRP_DB_FILE_DIR, osProxy.TempDir()).Return(errors.New("OsProxy.Setenv() failed"))
				fields.OsProxy = mockOsProxy
			},
			cleanup: func() {
				if err := testEnvSetter.UnsetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.UnsetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()
				tt.setup(mockCtrl, &tt.fields)
			}
			tr := NewTestEnvSetter(tt.fields.OsProxy)
			if err := tr.SetTestEnv(); (err != nil) != tt.wantErr {
				t.Errorf("TestEnvSetter.SetTestEnv() : got =\n%v, want =\n%v", err, tt.wantErr)
			}
			testWnJpnDBFileDirPath, err := dbFilePathProvider.GetWNJpnDBFileDirPath()
			if err != nil {
				t.Errorf("DBFilePathProvider.GetWNJpnDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
			}
			testJrpDBFileDirPath, err := dbFilePathProvider.GetJrpDBFileDirPath()
			if err != nil {
				t.Errorf("DBFilePathProvider.GetJrpDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
			}
			if testWnJpnDBFileDirPath != tt.wantWnJpnDBFileDirPath {
				t.Errorf("TestEnvSetter.SetTestEnv() : gotWnJpnDBFileDirPath =\n%v, wantWnJpnDBFileDirPath =\n%v", testWnJpnDBFileDirPath, tt.wantWnJpnDBFileDirPath)
			}
			if testJrpDBFileDirPath != tt.wantJrpDBFileDirPath {
				t.Errorf("TestEnvSetter.SetTestEnv() : gotJrpDBFileDirPath =\n%v, wantJrpDBFileDirPath =\n%v", testJrpDBFileDirPath, tt.wantJrpDBFileDirPath)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestTestEnvSetter_UnsetTestEnv(t *testing.T) {
	osProxy := osproxy.New()
	dbFilePathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	testEnvSetter := NewTestEnvSetter(osProxy)
	if err := testEnvSetter.UnsetTestEnv(); err != nil {
		t.Errorf("TestEnvSetter.UnsetTestEnv() : got =\n%v, want =\n%v", err, nil)
	}
	defaultWnJpnDBFileDirPath, err := dbFilePathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFilePathProvider.GetWNJpnDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
	}
	defaultJrpDBFileDirPath, err := dbFilePathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFilePathProvider.GetJrpDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
	}

	type fields struct {
		OsProxy osproxy.Os
	}
	tests := []struct {
		name                   string
		fields                 fields
		wantWnJpnDBFileDirPath string
		wantJrpDBFileDirPath   string
		wantErr                bool
		setup                  func(mockCtrl *gomock.Controller, tt *fields)
		cleanup                func()
	}{
		{
			name: "positive testing",
			fields: fields{
				OsProxy: osproxy.New(),
			},
			wantWnJpnDBFileDirPath: defaultWnJpnDBFileDirPath,
			wantJrpDBFileDirPath:   defaultJrpDBFileDirPath,
			wantErr:                false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := testEnvSetter.SetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.SetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
			cleanup: nil,
		},
		{
			name: "negative testing (OsProxy.Unsetenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR) failed)",
			fields: fields{
				OsProxy: nil,
			},
			wantWnJpnDBFileDirPath: osProxy.TempDir(),
			wantJrpDBFileDirPath:   osProxy.TempDir(),
			wantErr:                true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().Unsetenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR).Return(errors.New("OsProxy.Unsetenv() failed"))
				fields.OsProxy = mockOsProxy
				if err := testEnvSetter.SetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.SetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
			cleanup: func() {
				if err := testEnvSetter.UnsetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.UnsetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
		},
		{
			name: "negative testing (OsProxy.Unsetenv(dbfiledirpathprovider.JRP_ENV_JRP_DB_FILE_DIR) failed)",
			fields: fields{
				OsProxy: nil,
			},
			wantWnJpnDBFileDirPath: osProxy.TempDir(),
			wantJrpDBFileDirPath:   osProxy.TempDir(),
			wantErr:                true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().Unsetenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR).Return(nil)
				mockOsProxy.EXPECT().Unsetenv(dbfiledirpathprovider.JRP_ENV_JRP_DB_FILE_DIR).Return(errors.New("OsProxy.Unsetenv() failed"))
				fields.OsProxy = mockOsProxy
				if err := testEnvSetter.SetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.SetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
			cleanup: func() {
				if err := testEnvSetter.UnsetTestEnv(); err != nil {
					t.Errorf("TestEnvSetter.UnsetTestEnv() : got =\n%v, want =\n%v", err, nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()
				tt.setup(mockCtrl, &tt.fields)
			}
			tr := NewTestEnvSetter(tt.fields.OsProxy)
			if err := tr.UnsetTestEnv(); (err != nil) != tt.wantErr {
				t.Errorf("TestEnvSetter.UnsetTestEnv() : got =\n%v, want =\n%v", err, tt.wantErr)
			}
			testWnJpnDBFileDirPath, err := dbFilePathProvider.GetWNJpnDBFileDirPath()
			if err != nil {
				t.Errorf("DBFilePathProvider.GetWNJpnDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
			}
			testJrpDBFileDirPath, err := dbFilePathProvider.GetJrpDBFileDirPath()
			if err != nil {
				t.Errorf("DBFilePathProvider.GetJrpDBFileDirPath() : got =\n%v, want =\n%v", err, nil)
			}
			if testWnJpnDBFileDirPath != tt.wantWnJpnDBFileDirPath {
				t.Errorf("TestEnvSetter.UnsetTestEnv() : gotWnJpnDBFileDirPath =\n%v, wantWnJpnDBFileDirPath =\n%v", testWnJpnDBFileDirPath, tt.wantWnJpnDBFileDirPath)
			}
			if testJrpDBFileDirPath != tt.wantJrpDBFileDirPath {
				t.Errorf("TestEnvSetter.UnsetTestEnv() : gotJrpDBFileDirPath =\n%v, wantJrpDBFileDirPath =\n%v", testJrpDBFileDirPath, tt.wantJrpDBFileDirPath)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
