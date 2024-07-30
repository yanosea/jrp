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
		downloader *DBFileDownloader
		env        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name:    "positive testing (no env)",
			args:    args{downloader: nil, env: ""},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, OSFileSystem{}, DefaultHttpClient{}, DefaultIO{}, DefaultGzip{})
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "positive testing (with env)",
			args:    args{downloader: nil, env: filepath.Join(tcu.HomeDir, "jrp")},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, OSFileSystem{}, DefaultHttpClient{}, DefaultIO{}, DefaultGzip{})
				dbFilePath := filepath.Join(tt.env, constant.WNJPN_DB_FILE_NAME)
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (user.Current() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mu := mock_logic.NewMockUser(mockCtrl)
				mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))
				tt.downloader = NewDBFileDownloader(OsEnv{}, mu, OSFileSystem{}, DefaultHttpClient{}, DefaultIO{}, DefaultGzip{})
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (http.Get() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mhc := mock_logic.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("failed to get http response"))
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, OSFileSystem{}, mhc, DefaultIO{}, DefaultGzip{})
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (os.Create(tempFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				mfs.EXPECT().Create(gomock.Any()).Return(nil, errors.New("failed to create temp file"))
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, mfs, DefaultHttpClient{}, DefaultIO{}, DefaultGzip{})
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (io.Copy(out, resp.Body) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_logic.NewMockIO(mockCtrl)
				mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy response body to temp file"))
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, OSFileSystem{}, DefaultHttpClient{}, mio, DefaultGzip{})
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (gzip.NewReader() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mgz := mock_logic.NewMockGzip(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(nil, errors.New("failed to new reader"))
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, OSFileSystem{}, DefaultHttpClient{}, DefaultIO{}, mgz)
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (os.Create(dbFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
				tempFile, _ := os.Create(tempFilePath)
				mfs.EXPECT().Create(tempFilePath).Return(tempFile, nil)
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp", constant.WNJPN_DB_FILE_NAME)
				mfs.EXPECT().Create(dbFilePath).Return(nil, errors.New("failed to create db file"))
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, mfs, DefaultHttpClient{}, DefaultIO{}, DefaultGzip{})
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (io.Copy(f, gz) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_logic.NewMockIO(mockCtrl)
				gomock.InOrder(
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1024), nil),
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy decompressed file to db file")),
				)
				mfs := mock_logic.NewMockFileSystem(mockCtrl)
				tempFile, err := os.CreateTemp("", "test")
				if err != nil {
					t.Fatal(err)
				}
				tempFile.Close()
				mfs.EXPECT().Create(gomock.Any()).Return(tempFile, nil).AnyTimes()
				mhc := mock_logic.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte{}))}, nil).AnyTimes()
				mgz := mock_logic.NewMockGzip(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte{})), nil).AnyTimes()
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, mfs, mhc, mio, mgz)
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFilePath)
			},
		}, {
			name:    "negative testing (os.Remove(tempFilePath) fails)",
			args:    args{downloader: nil, env: ""},
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
				tt.downloader = NewDBFileDownloader(OsEnv{}, OsUser{}, mfs, DefaultHttpClient{}, DefaultIO{}, DefaultGzip{})
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

			if err := tt.args.downloader.Download(); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
