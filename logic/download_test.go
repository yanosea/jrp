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

	type args struct {
		e   Env
		u   User
		fs  FileSystem
		hc  HttpClient
		io  IO
		gz  Gzip
		env string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name:    "positive testing (no env)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: false,
			setup:   nil,
		}, {
			name:    "positive testing (with env)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: filepath.Join(tcu.HomeDir, "jrp")},
			wantErr: false,
			setup:   nil,
		}, {
			name:    "negative testing (user.Current() fails)",
			args:    args{e: OsEnv{}, u: nil, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mu := mock_logic.NewMockUser(mockCtrl)
				mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))
				tt.u = mu
			},
		}, {
			name:    "negative testing (http.Get() fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: nil, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mhc := mock_logic.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("failed to get http response"))
				tt.hc = mhc
			},
		}, {
			name:    "negative testing (os.Create(tempFilePath) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: nil, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				mfs.EXPECT().Create(gomock.Any()).Return(nil, errors.New("failed to create temp file"))
				tt.fs = mfs
			},
		}, {
			name:    "negative testing (io.Copy(out, resp.Body) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: nil, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_logic.NewMockIO(mockCtrl)
				mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy response body to temp file"))
				tt.io = mio
			},
		}, {
			name:    "negative testing (gzip.NewReader() fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mgz := mock_logic.NewMockGzip(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(nil, errors.New("failed to new reader"))
				tt.gz = mgz
			},
		}, {
			name:    "negative testing (os.Create(dbFilePath) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: nil, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: filepath.Join(tcu.HomeDir, "jrp")},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
				tempFile, _ := os.Open(tempFilePath)
				mfs.EXPECT().Create(tempFilePath).Return(tempFile, nil)
				dbFilePath := filepath.Join(tt.env, constant.WNJPN_DB_FILE_NAME)
				mfs.EXPECT().Create(dbFilePath).Return(nil, errors.New("failed to create db file"))
				tt.fs = mfs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.setup != nil {
				tt.setup(ctrl, &tt.args)
			}

			dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
			if tt.args.env != "" {
				os.Setenv(constant.JRP_ENV_WORDNETJP_DIR, tt.args.env)
				defer os.Unsetenv(constant.JRP_ENV_WORDNETJP_DIR)
				dbFilePath = os.Getenv(constant.JRP_ENV_WORDNETJP_DIR)
			}
			os.RemoveAll(dbFilePath)
			if err := Download(tt.args.e, tt.args.u, tt.args.fs, tt.args.hc, tt.args.io, tt.args.gz); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
