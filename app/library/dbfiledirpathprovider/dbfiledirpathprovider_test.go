package dbfiledirpathprovider

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/user"

	"github.com/yanosea/jrp/mock/app/proxy/user"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	userProxy := userproxy.New()

	type args struct {
		filepath filepathproxy.FilePath
		os       osproxy.Os
		user     userproxy.User
	}
	tests := []struct {
		name string
		args args
		want *DBFileDirPathProvider
	}{
		{
			name: "positive testing",
			args: args{
				filepath: filepathProxy,
				os:       osProxy,
				user:     userProxy,
			},
			want: &DBFileDirPathProvider{
				FilepathProxy: filepathProxy,
				OsProxy:       osProxy,
				UserProxy:     userProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.filepath, tt.args.os, tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want=\n%v", got, tt.want)
			}
		})
	}
}

func TestDBFileDirPathProvider_GetJrpDBFileDirPath(t *testing.T) {
	userProxy := userproxy.New()
	currentUser, err := userProxy.Current()
	if err != nil {
		t.Errorf("UserProxy.Current() : error =\n%v", err)
	}

	type fields struct {
		FilepathProxy filepathproxy.FilePath
		OsProxy       osproxy.Os
		UserProxy     userproxy.User
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (with no xdg data home, with no env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/.local/share/jrp"),
			wantErr: false,
			setup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_DB_FILE_DIR")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_DB_FILE_DIR")
			},
		},
		{
			name: "positive testing (with xdg data home, with no env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/.local/share/jrp"),
			wantErr: false,
			setup: func() {
				os.Setenv("XDG_DATA_HOME", currentUser.FieldUser.HomeDir+"/.local/share")
				os.Unsetenv("JRP_DB_FILE_DIR")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_DB_FILE_DIR")
			},
		},
		{
			name: "positive testing (with no xdg data home, with env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/jrp"),
			wantErr: false,
			setup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Setenv("JRP_DB_FILE_DIR", currentUser.FieldUser.HomeDir+"/jrp")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_DB_FILE_DIR")
			},
		},
		{
			name: "positive testing (with xdg data home, with env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/jrp"),
			wantErr: false,
			setup: func() {
				os.Setenv("XDG_DATA_HOME", currentUser.FieldUser.HomeDir+"/.local/share")
				os.Setenv("JRP_DB_FILE_DIR", currentUser.FieldUser.HomeDir+"/jrp")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_DB_FILE_DIR")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			d := New(
				tt.fields.FilepathProxy,
				tt.fields.OsProxy,
				tt.fields.UserProxy,
			)
			got, err := d.GetJrpDBFileDirPath()
			gotDir := testutility.ReplaceDoubleSlashToSingleSlash(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if gotDir != tt.want {
				t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestDBFileDirPathProvider_GetWNJpnDBFileDirPath(t *testing.T) {
	userProxy := userproxy.New()
	currentUser, err := userProxy.Current()
	if err != nil {
		t.Errorf("UserProxy.Current() : error =\n%v", err)
	}

	type fields struct {
		FilepathProxy filepathproxy.FilePath
		OsProxy       osproxy.Os
		UserProxy     userproxy.User
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (with no xdg data home, with no env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/.local/share/jrp"),
			wantErr: false,
			setup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_WNJPN_DB_FILE_DIR")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_WNJPN_DB_FILE_DIR")
			},
		},
		{
			name: "positive testing (with xdg data home, with no env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/.local/share/jrp"),
			wantErr: false,
			setup: func() {
				os.Setenv("XDG_DATA_HOME", currentUser.FieldUser.HomeDir+"/.local/share")
				os.Unsetenv("JRP_WNJPN_DB_FILE_DIR")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_WNJPN_DB_FILE_DIR")
			},
		},
		{
			name: "positive testing (with no xdg data home, with env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/jrp"),
			wantErr: false,
			setup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Setenv("JRP_WNJPN_DB_FILE_DIR", currentUser.FieldUser.HomeDir+"/jrp")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_WNJPN_DB_FILE_DIR")
			},
		},
		{
			name: "positive testing (with xdg data home, with env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/jrp"),
			wantErr: false,
			setup: func() {
				os.Setenv("XDG_DATA_HOME", currentUser.FieldUser.HomeDir+"/.local/share")
				os.Setenv("JRP_WNJPN_DB_FILE_DIR", currentUser.FieldUser.HomeDir+"/jrp")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("JRP_WNJPN_DB_FILE_DIR")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			d := New(
				tt.fields.FilepathProxy,
				tt.fields.OsProxy,
				tt.fields.UserProxy,
			)
			got, err := d.GetWNJpnDBFileDirPath()
			gotDir := testutility.ReplaceDoubleSlashToSingleSlash(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if gotDir != tt.want {
				t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestDBFileDirPathProvider_getDBFileDirPath(t *testing.T) {
	userProxy := userproxy.New()
	currentUser, err := userProxy.Current()
	if err != nil {
		t.Errorf("UserProxy.Current() : error =\n%v", err)
	}

	type fields struct {
		FilepathProxy filepathproxy.FilePath
		OsProxy       osproxy.Os
		UserProxy     userproxy.User
	}
	type args struct {
		envVar string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (with no xdg data home, with no env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			args: args{
				envVar: "TEST_ENV",
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/.local/share/jrp"),
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("TEST_ENV")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("TEST_ENV")
			},
		},
		{
			name: "positive testing (with xdg data home, with no env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			args: args{
				envVar: "TEST_ENV",
			},
			want:    testutility.ReplaceDoubleSlashToSingleSlash(currentUser.FieldUser.HomeDir + "/.local/share/jrp"),
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				os.Setenv("XDG_DATA_HOME", currentUser.FieldUser.HomeDir+"/.local/share")
				os.Unsetenv("TEST_ENV")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("TEST_ENV")
			},
		},
		{
			name: "positive testing (with no xdg data home, with env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			args: args{
				envVar: "TEST_ENV",
			},
			want:    "TEST_DIR",
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				os.Unsetenv("XDG_DATA_HOME")
				os.Setenv("TEST_ENV", "TEST_DIR")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("TEST_ENV")
			},
		},
		{
			name: "positive testing (with xdg data home, with env)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     userproxy.New(),
			},
			args: args{
				envVar: "TEST_ENV",
			},
			want:    "TEST_DIR",
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				os.Setenv("XDG_DATA_HOME", currentUser.FieldUser.HomeDir+"/.local/share")
				os.Setenv("TEST_ENV", "TEST_DIR")
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
				os.Unsetenv("TEST_ENV")
			},
		},
		{
			name: "negative testing (UserProxy.Current() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				UserProxy:     nil,
			},
			args: args{
				envVar: "TEST_ENV",
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				os.Unsetenv("TEST_ENV")
				mockUserProxy := mockuserproxy.NewMockUser(mockCtrl)
				mockUserProxy.EXPECT().Current().Return(nil, errors.New("UserProxy.Current() failed"))
				tt.UserProxy = mockUserProxy
			},
			cleanup: func() {
				os.Unsetenv("TEST_ENV")
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
			d := New(
				tt.fields.FilepathProxy,
				tt.fields.OsProxy,
				tt.fields.UserProxy,
			)
			got, err := d.getDBFileDirPath(tt.args.envVar)
			gotDir := testutility.ReplaceDoubleSlashToSingleSlash(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBFileDirPathProvider.getDBFileDirPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDir != tt.want {
				t.Errorf("DBFileDirPathProvider.getDBFileDirPath() = %v, want %v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
