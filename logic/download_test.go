package logic

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/constant"
	mock_logic "github.com/yanosea/jrp/mock/logic"
)

func TestDownload(t *testing.T) {
	tu := OsUser{}
	tcu, _ := tu.Current()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mu := mock_logic.NewMockUser(ctrl)
	mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))

	mhc := mock_logic.NewMockHTTPClient(ctrl)
	mhc.EXPECT().Get(constant.WNJPN_DB_ARCHIVE_FILE_URL).Return(nil, errors.New("failed to get http response"))

	type args struct {
		e   Env
		u   User
		fs  FileSystem
		hc  HttpClient
		env string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive testing (no env)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, env: ""},
			wantErr: false,
		}, {
			name:    "positive testing (with env)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, env: filepath.Join(tcu.HomeDir, "jrp")},
			wantErr: false,
		}, {
			name:    "negative testing (user.Current() fails)",
			args:    args{e: OsEnv{}, u: mu, fs: OSFileSystem{}, hc: DefaultHttpClient{}, env: ""},
			wantErr: true,
		}, {
			name:    "negative testing (http.Get() fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: mhc, env: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
			if tt.args.env != "" {
				os.Setenv(constant.JRP_ENV_WORDNETJP_DIR, tt.args.env)
				defer os.Unsetenv(constant.JRP_ENV_WORDNETJP_DIR)
				dbFilePath = os.Getenv(constant.JRP_ENV_WORDNETJP_DIR)
			}
			os.RemoveAll(dbFilePath)
			if err := Download(tt.args.e, tt.args.u, tt.args.fs, tt.args.hc); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
