package logic

import (
	"bytes"
	"errors"
	"io"
	"net/http"
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
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "positive testing (with env)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: filepath.Join(tcu.HomeDir, "jrp")},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				dbFilePath := filepath.Join(tt.env, constant.WNJPN_DB_FILE_NAME)
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (user.Current() fails)",
			args:    args{e: OsEnv{}, u: nil, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mu := mock_logic.NewMockUser(mockCtrl)
				mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))
				tt.u = mu
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (http.Get() fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: nil, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mhc := mock_logic.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("failed to get http response"))
				tt.hc = mhc
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (os.Create(tempFilePath) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: nil, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				mfs.EXPECT().Create(gomock.Any()).Return(nil, errors.New("failed to create temp file"))
				tt.fs = mfs
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (io.Copy(out, resp.Body) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: nil, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_logic.NewMockIO(mockCtrl)
				mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy response body to temp file"))
				tt.io = mio
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (gzip.NewReader() fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mgz := mock_logic.NewMockGzip(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(nil, errors.New("failed to new reader"))
				tt.gz = mgz
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (os.Create(dbFilePath) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: nil, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
				tempFile, _ := os.Create(tempFilePath)
				mfs.EXPECT().Create(tempFilePath).Return(tempFile, nil)
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp", constant.WNJPN_DB_FILE_NAME)
				mfs.EXPECT().Create(dbFilePath).Return(nil, errors.New("failed to create db file"))
				tt.fs = mfs
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (io.Copy(f, gz) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: OSFileSystem{}, hc: DefaultHttpClient{}, io: nil, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_logic.NewMockIO(mockCtrl)
				gomock.InOrder(
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1024), nil),
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy decompressed file to db file")),
				)
				tt.io = mio
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				tempFile, err := os.CreateTemp("", "test")
				if err != nil {
					t.Fatal(err)
				}
				tempFile.Close()
				mfs.EXPECT().Create(gomock.Any()).Return(tempFile, nil).AnyTimes()
				tt.fs = mfs
				mhc := mock_logic.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte{}))}, nil).AnyTimes()
				tt.hc = mhc
				mgz := mock_logic.NewMockGzip(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte{})), nil).AnyTimes()
				tt.gz = mgz
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (os.Remove(tempFilePath) fails)",
			args:    args{e: OsEnv{}, u: OsUser{}, fs: nil, hc: DefaultHttpClient{}, io: DefaultIO{}, gz: DefaultGzip{}, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
				tempFile, _ := os.Create(tempFilePath)
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp", constant.WNJPN_DB_FILE_NAME)
				dbFile, _ := os.Create(dbFilePath)
				gomock.InOrder(
					mfs.EXPECT().Create(tempFilePath).Return(tempFile, nil),
					mfs.EXPECT().Create(dbFilePath).Return(dbFile, nil),
					mfs.EXPECT().Remove(gomock.Any()).Return(errors.New("failed to remove temp file")),
				)
				tt.fs = mfs
				os.RemoveAll(dbFilePath)
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

			if err := Download(tt.args.e, tt.args.u, tt.args.fs, tt.args.hc, tt.args.io, tt.args.gz); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
