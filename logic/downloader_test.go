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
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/usermanager"
	mock_fs "github.com/yanosea/jrp/mock/fs"
	mock_gzip "github.com/yanosea/jrp/mock/gzip"
	mock_httpclient "github.com/yanosea/jrp/mock/httpclient"
	mock_iomanager "github.com/yanosea/jrp/mock/iomanager"
	mock_usermanager "github.com/yanosea/jrp/mock/usermanager"
)

func TestDownload(t *testing.T) {
	tu := usermanager.OSUserProvider{}
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
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
			},
		}, {
			name:    "positive testing (with env)",
			args:    args{downloader: nil, env: filepath.Join(tcu.HomeDir, "jrp")},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				dbFileDirPath := filepath.Join(tt.env)
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (user.Current() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mu := mock_usermanager.NewMockUserProvider(mockCtrl)
				mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(mu, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (http.Get() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mhc := mock_httpclient.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("failed to get http response"))
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, mhc, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (os.Create(tempFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				mfs.EXPECT().Create(gomock.Any()).Return(nil, errors.New("failed to create temp file"))
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (io.Copy(out, resp.Body) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_iomanager.NewMockIOHelper(mockCtrl)
				mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy response body to temp file"))
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, mio, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (gzip.NewReader() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mgz := mock_gzip.NewMockGzipHandler(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(nil, errors.New("failed to new reader"))
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, mgz)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (os.Create(dbFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
				tempFile, _ := os.Create(tempFilePath)
				mfs.EXPECT().Create(tempFilePath).Return(tempFile, nil)
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp", constant.WNJPN_DB_FILE_NAME)
				mfs.EXPECT().Create(dbFilePath).Return(nil, errors.New("failed to create db file"))
				os.RemoveAll(dbFilePath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (io.Copy(f, gz) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_iomanager.NewMockIOHelper(mockCtrl)
				gomock.InOrder(
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1024), nil),
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy decompressed file to db file")),
				)
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				tempFile, err := os.CreateTemp("", "test")
				if err != nil {
					t.Fatal(err)
				}
				tempFile.Close()
				mfs.EXPECT().Create(gomock.Any()).Return(tempFile, nil).AnyTimes()
				mhc := mock_httpclient.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte{}))}, nil).AnyTimes()
				mgz := mock_gzip.NewMockGzipHandler(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte{})), nil).AnyTimes()
				dbFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
				os.RemoveAll(dbFileDirPath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, mhc, mio, mgz)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (os.Remove(tempFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
				tempFile, _ := os.Create(tempFilePath)
				dbFilePath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp", constant.WNJPN_DB_FILE_NAME)
				dbFile, _ := os.Create(dbFilePath)
				gomock.InOrder(
					mfs.EXPECT().Create(tempFilePath).Return(tempFile, nil),
					mfs.EXPECT().Create(dbFilePath).Return(dbFile, nil),
					mfs.EXPECT().Remove(gomock.Any()).Return(errors.New("failed to remove temp file")),
				)
				os.RemoveAll(dbFilePath)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
				tt.downloader = downloader
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

	defaultDBFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	os.RemoveAll(defaultDBFileDirPath)
	envDBFileDirPath := filepath.Join(tcu.HomeDir, "jrp")
	os.RemoveAll(envDBFileDirPath)
}
