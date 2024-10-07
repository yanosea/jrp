package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/downloader"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/gzip"
	"github.com/yanosea/jrp/app/proxy/http"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/spinner"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"

	"github.com/yanosea/jrp/mock/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/mock/app/library/utility"
	"github.com/yanosea/jrp/mock/app/proxy/spinner"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNewDownloadCommand(t *testing.T) {
	type args struct {
		g *GlobalOption
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name: "positive testing",
			args: args{
				g: NewGlobalOption(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDownloadCommand(tt.args.g)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewDownloadCommand().Execute() : error =\n%v", err)
			}
		})
	}
}

func Test_downloadOption_downloadRunE(t *testing.T) {
	osProxy := osproxy.New()
	globalOption := NewGlobalOption(fmtproxy.New(), osProxy, strconvproxy.New())
	filepathProxy := filepathproxy.New()
	dbFilePathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFilePathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFilePathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Downloader            downloader.Downloadable
		SpinnerProxy          spinnerproxy.Spinner
		Utility               utility.UtilityInterface
	}
	type args struct {
		in0 *cobra.Command
		in1 []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:                   globalOption.Out,
				ErrOut:                globalOption.ErrOut,
				DBFileDirPathProvider: dbFilePathProvider,
				Downloader: downloader.New(
					filepathproxy.New(),
					gzipproxy.New(),
					httpproxy.New(),
					ioproxy.New(),
					osproxy.New(),
					globalOption.Utility,
				),
				SpinnerProxy: spinnerproxy.New(),
				Utility:      globalOption.Utility,
			},
			args: args{
				in0: nil,
				in1: nil,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (DBFilePathProvider.GetWNJpnDBFileDirPath() failed)",
			fields: fields{
				Out:                   globalOption.Out,
				ErrOut:                globalOption.ErrOut,
				DBFileDirPathProvider: nil,
				Downloader: downloader.New(
					filepathproxy.New(),
					gzipproxy.New(),
					httpproxy.New(),
					ioproxy.New(),
					osproxy.New(),
					globalOption.Utility,
				),
				SpinnerProxy: spinnerproxy.New(),
				Utility:      globalOption.Utility,
			},
			args: args{
				in0: nil,
				in1: nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockDBFileDirPathProvider := mockdbfiledirpathprovider.NewMockDBFileDirPathProvidable(mockCtrl)
				mockDBFileDirPathProvider.EXPECT().GetWNJpnDBFileDirPath().Return("", errors.New("DBFileDirPathProvider.GetWNJpnDBFileDirPath() failed"))
				tt.DBFileDirPathProvider = mockDBFileDirPathProvider
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (Utility.CreateDirIfNotExist() failed)",
			fields: fields{
				Out:                   globalOption.Out,
				ErrOut:                globalOption.ErrOut,
				DBFileDirPathProvider: dbFilePathProvider,
				Downloader: downloader.New(
					filepathproxy.New(),
					gzipproxy.New(),
					httpproxy.New(),
					ioproxy.New(),
					osproxy.New(),
					globalOption.Utility,
				),
				SpinnerProxy: spinnerproxy.New(),
				Utility:      nil,
			},
			args: args{
				in0: nil,
				in1: nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockUtility := mockutility.NewMockUtilityInterface(mockCtrl)
				mockUtility.EXPECT().CreateDirIfNotExist(gomock.Any()).Return(errors.New("Utility.CreateDirIfNotExist() failed"))
				tt.Utility = mockUtility
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
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
			o := &downloadOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Downloader:            tt.fields.Downloader,
				SpinnerProxy:          tt.fields.SpinnerProxy,
				Utility:               tt.fields.Utility,
			}
			if err := o.downloadRunE(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("downloadOption.downloadRunE() : error =\n%v, wantErr \n%v", err, tt.wantErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_downloadOption_download(t *testing.T) {
	osProxy := osproxy.New()
	globalOption := NewGlobalOption(fmtproxy.New(), osProxy, strconvproxy.New())
	filepathProxy := filepathproxy.New()
	dbFilePathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFilePathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFilePathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Downloader            downloader.Downloadable
		SpinnerProxy          spinnerproxy.Spinner
		Utility               utility.UtilityInterface
	}
	type args struct {
		wnJpnDBFileDirPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:                   globalOption.Out,
				ErrOut:                globalOption.ErrOut,
				DBFileDirPathProvider: dbFilePathProvider,
				Downloader: downloader.New(
					filepathproxy.New(),
					gzipproxy.New(),
					httpproxy.New(),
					ioproxy.New(),
					osproxy.New(),
					globalOption.Utility,
				),
				SpinnerProxy: spinnerproxy.New(),
				Utility:      globalOption.Utility,
			},
			args: args{
				wnJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SpinnerProxy.SetColor() failed)",
			fields: fields{
				Out:                   globalOption.Out,
				ErrOut:                globalOption.ErrOut,
				DBFileDirPathProvider: dbFilePathProvider,
				Downloader: downloader.New(
					filepathproxy.New(),
					gzipproxy.New(),
					httpproxy.New(),
					ioproxy.New(),
					osproxy.New(),
					globalOption.Utility,
				),
				SpinnerProxy: nil,
				Utility:      globalOption.Utility,
			},
			args: args{
				wnJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockSpinnerInstance := mockspinnerproxy.NewMockSpinnerInstanceInterface(mockCtrl)
				mockSpinnerInstance.EXPECT().Reverse()
				mockSpinnerInstance.EXPECT().SetColor(gomock.Any()).Return(errors.New("SpinnerProxy.SetColor() failed"))
				mockSpinnerProxy := mockspinnerproxy.NewMockSpinner(mockCtrl)
				mockSpinnerProxy.EXPECT().NewSpinner().Return(mockSpinnerInstance)
				tt.SpinnerProxy = mockSpinnerProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
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
			o := &downloadOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Downloader:            tt.fields.Downloader,
				SpinnerProxy:          tt.fields.SpinnerProxy,
				Utility:               tt.fields.Utility,
			}
			if err := o.download(tt.args.wnJpnDBFileDirPath); (err != nil) != tt.wantErr {
				t.Errorf("downloadOption.download() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_downloadOption_writeDownloadResult(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	globalOption := NewGlobalOption(fmtproxy.New(), osProxy, strconvproxy.New())
	globalOption.Out = capturer.OutBuffer
	globalOption.ErrOut = capturer.ErrBuffer
	o := &downloadOption{
		Out:     globalOption.Out,
		ErrOut:  globalOption.ErrOut,
		Utility: globalOption.Utility,
	}
	colorProxy := colorproxy.New()

	type fields struct {
		t        *testing.T
		fnc      func()
		capturer *testutility.Capturer
	}
	tests := []struct {
		name       string
		fields     fields
		wantStdOut string
		wantStdErr string
		wantErr    bool
	}{
		{
			name: "positive testing (result is DownloadedFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					o.writeDownloadResult(downloader.DownloadedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.DOWNLOAD_MESSAGE_FAILED) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (result is DownloadedAlready)",
			fields: fields{
				t: t,
				fnc: func() {
					o.writeDownloadResult(downloader.DownloadedAlready)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is DownloadedSuccessfully)",
			fields: fields{
				t: t,
				fnc: func() {
					o.writeDownloadResult(downloader.DownloadedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.DOWNLOAD_MESSAGE_SUCCEEDED) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if stdout != tt.wantStdOut {
				t.Errorf("DownloadOption.writeDownloadResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("DownloadOption.writeDownloadResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
