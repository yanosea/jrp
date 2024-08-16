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
	"github.com/yanosea/jrp/internal/spinnerservice"
	"github.com/yanosea/jrp/internal/usermanager"

	mock_fs "github.com/yanosea/jrp/mock/fs"
	mock_gzip "github.com/yanosea/jrp/mock/gzip"
	mock_httpclient "github.com/yanosea/jrp/mock/httpclient"
	mock_iomanager "github.com/yanosea/jrp/mock/iomanager"
	mock_spinnerservice "github.com/yanosea/jrp/mock/spinnerservice"
	mock_usermanager "github.com/yanosea/jrp/mock/usermanager"
)

func TestNewDBFileDownloader(t *testing.T) {
	type args struct {
		u usermanager.UserProvider
		f fs.FileManager
		h httpclient.HTTPClient
		i iomanager.IOHelper
		g gzip.GzipHandler
		s spinnerservice.SpinnerService
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{u: usermanager.OSUserProvider{}, f: fs.OsFileManager{}, h: httpclient.DefaultHTTPClient{}, i: iomanager.DefaultIOHelper{}, g: gzip.DefaultGzipHandler{}, s: spinnerservice.NewRealSpinnerService()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewDBFileDownloader(tt.args.u, tt.args.f, tt.args.h, tt.args.i, tt.args.g, tt.args.s)
			if u == nil {
				t.Errorf("NewDBFileDownloader() : returned nil")
			}
		})
	}
}

func TestDownload(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()
	sp := spinnerservice.NewRealSpinnerService()
	defaultDBFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	envDBFileDirPath := filepath.Join(tcu.HomeDir, "jrp")
	defaultTempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)

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
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "positive testing (with env)",
			args:    args{downloader: nil, env: filepath.Join(tcu.HomeDir, "jrp")},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "positive testing (already has db file)",
			args:    args{downloader: nil, env: ""},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
				if err := downloader.Download(); err != nil {
					t.Error(err)
				}
			},
		}, {
			name:    "negative testing (user.Current() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mu := mock_usermanager.NewMockUserProvider(mockCtrl)
				mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))
				downloader := NewDBFileDownloader(mu, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (os.MkdirAll(dbFileDirPath, os.Filemode(0755) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				mfs.EXPECT().MkdirAll(gomock.Any(), os.FileMode(0755)).Return(errors.New("failed to create db file dir"))
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (d.Spinner.SetColor(\"yellow\") fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				ms := mock_spinnerservice.NewMockSpinnerService(mockCtrl)
				ms.EXPECT().SetColor(gomock.Any()).Return(errors.New("failed to set spinner color"))
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, ms)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (http.Get() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mhc := mock_httpclient.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("failed to get http response"))
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, mhc, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (os.Create(tempFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				gomock.InOrder(
					mfs.EXPECT().MkdirAll(defaultDBFileDirPath, os.FileMode(0755)).Return(nil),
					mfs.EXPECT().Create(gomock.Any()).Return(nil, errors.New("failed to create temp file")),
				)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (io.Copy(out, resp.Body) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mio := mock_iomanager.NewMockIOHelper(mockCtrl)
				mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy response body to temp file"))
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, mio, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (out.Seek(0, 0) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mf := mock_fs.NewMockFile(mockCtrl)
				gomock.InOrder(
					mf.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), errors.New("failed to seek to start")),
					mf.EXPECT().Close().Return(nil),
				)
				var mfAsFsFile fs.File = mf
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				gomock.InOrder(
					mfs.EXPECT().MkdirAll(defaultDBFileDirPath, os.FileMode(0755)).Return(nil),
					mfs.EXPECT().Create(defaultTempFilePath).Return(mfAsFsFile, nil),
				)
				mio := mock_iomanager.NewMockIOHelper(mockCtrl)
				mio.EXPECT().Copy(mf, gomock.Any()).Return(int64(0), nil)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, httpclient.DefaultHTTPClient{}, mio, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (gzip.NewReader() fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mgz := mock_gzip.NewMockGzipHandler(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(nil, errors.New("failed to new reader"))
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, mgz, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (os.Create(dbFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				tempFile, _ := os.Create(defaultTempFilePath)
				dbFilePath := filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME)
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				gomock.InOrder(
					mfs.EXPECT().MkdirAll(defaultDBFileDirPath, os.FileMode(0755)).Return(nil),
					mfs.EXPECT().Create(defaultTempFilePath).Return(tempFile, nil),
					mfs.EXPECT().Create(dbFilePath).Return(nil, errors.New("failed to create db file")),
				)
				downloader := NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, sp)
				tt.downloader = downloader
			},
		}, {
			name:    "negative testing (io.Copy(f, gz) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				tempFile, _ := os.Create(defaultTempFilePath)
				dbFilePath := filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME)
				dbFile, _ := os.Create(dbFilePath)
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				gomock.InOrder(
					mfs.EXPECT().MkdirAll(defaultDBFileDirPath, os.FileMode(0755)).Return(nil),
					mfs.EXPECT().Create(defaultTempFilePath).Return(tempFile, nil),
					mfs.EXPECT().Create(dbFilePath).Return(dbFile, nil),
				)
				mio := mock_iomanager.NewMockIOHelper(mockCtrl)
				gomock.InOrder(
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1024), nil),
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to copy decompressed file to db file")),
				)
				mhc := mock_httpclient.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte{}))}, nil)
				mgz := mock_gzip.NewMockGzipHandler(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte{})), nil)
				tt.downloader = NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, mhc, mio, mgz, sp)
			},
		}, {
			name:    "negative testing (os.Remove(tempFilePath) fails)",
			args:    args{downloader: nil, env: ""},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				tempFile, _ := os.Create(defaultTempFilePath)
				dbFilePath := filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME)
				dbFile, _ := os.Create(dbFilePath)
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				gomock.InOrder(
					mfs.EXPECT().MkdirAll(defaultDBFileDirPath, os.FileMode(0755)).Return(nil),
					mfs.EXPECT().Create(defaultTempFilePath).Return(tempFile, nil),
					mfs.EXPECT().Create(dbFilePath).Return(dbFile, nil),
					mfs.EXPECT().RemoveAll(gomock.Any()).Return(errors.New("failed to remove temp file")),
				)
				mio := mock_iomanager.NewMockIOHelper(mockCtrl)
				gomock.InOrder(
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1024), nil),
					mio.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1024), nil),
				)
				mhc := mock_httpclient.NewMockHTTPClient(mockCtrl)
				mhc.EXPECT().Get(gomock.Any()).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte{}))}, nil)
				mgz := mock_gzip.NewMockGzipHandler(mockCtrl)
				mgz.EXPECT().NewReader(gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte{})), nil)
				os.RemoveAll(filepath.Join(tcu.HomeDir, ".local", "share", "jrp"))
				tt.downloader = NewDBFileDownloader(usermanager.OSUserProvider{}, mfs, mhc, mio, mgz, sp)
			},
		},
	}
	for _, tt := range tests {
		os.RemoveAll(defaultDBFileDirPath)
		os.RemoveAll(envDBFileDirPath)
		os.RemoveAll(defaultTempFilePath)
		os.Setenv(constant.JRP_ENV_WORDNETJP_DIR, tt.args.env)
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.setup != nil {
				tt.setup(ctrl, &tt.args)
			}

			if err := tt.args.downloader.Download(); (err != nil) != tt.wantErr {
				t.Errorf("Download() : error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		os.RemoveAll(defaultDBFileDirPath)
		os.RemoveAll(envDBFileDirPath)
		os.RemoveAll(defaultTempFilePath)
	}
}
