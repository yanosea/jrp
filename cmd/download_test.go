package cmd_test

import (
	"errors"
	"os"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/cmd"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"

	mock_downloader "github.com/yanosea/jrp/mock/downloader"
)

func TestNewDownloadCommand(t *testing.T) {
	type args struct {
		globalOption *cmd.GlobalOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{globalOption: &cmd.GlobalOption{Out: os.Stdout, ErrOut: os.Stderr}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.NewDownloadCommand(tt.args.globalOption)
			if got == nil {
				t.Errorf("NewDownloadCommand() : returned nil")
			}
		})
	}
}

func TestDownloadRunE(t *testing.T) {
	type args struct {
		o cmd.DownloadOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{o: cmd.DownloadOption{Out: os.Stdout, ErrOut: os.Stderr, Downloader: logic.NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := tt.args.o.DownloadRunE(nil, nil)
			if err != nil {
				t.Errorf("DownloadRunE() : error = %v", err)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	type args struct {
		o cmd.DownloadOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name:    "positive testing",
			args:    args{o: cmd.DownloadOption{Out: os.Stdout, ErrOut: os.Stderr, Downloader: logic.NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})}},
			wantErr: false,
			setup:   nil,
		},
		{
			name:    "negative testing (Download() fails)",
			args:    args{o: cmd.DownloadOption{Out: os.Stdout, ErrOut: os.Stderr, Downloader: nil}},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				md := mock_downloader.NewMockDownloader(mockCtrl)
				md.EXPECT().Download().Return(errors.New("failed to download db file"))
				tt.o.Downloader = md
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

			err := tt.args.o.Download()
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() : error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
