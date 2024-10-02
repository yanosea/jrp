package downloader

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/gzip"
	"github.com/yanosea/jrp/app/proxy/http"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/user"

	"github.com/yanosea/jrp/mock/app/library/utility"
	"github.com/yanosea/jrp/mock/app/proxy/gzip"
	"github.com/yanosea/jrp/mock/app/proxy/http"
	"github.com/yanosea/jrp/mock/app/proxy/io"
	"github.com/yanosea/jrp/mock/app/proxy/os"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	filepathProxy := filepathproxy.New()
	gzipProxy := gzipproxy.New()
	httpProxy := httpproxy.New()
	ioProxy := ioproxy.New()
	osProxy := osproxy.New()
	util := utility.New(
		fmtproxy.New(),
		osProxy,
		strconvproxy.New(),
	)
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}

	type args struct {
		filepathProxy      filepathproxy.FilePath
		gzipProxy          gzipproxy.Gzip
		httpProxy          httpproxy.Http
		ioProxy            ioproxy.Io
		osProxy            osproxy.Os
		utility            utility.UtilityInterface
		wnJpnDBFileDirPath string
	}
	tests := []struct {
		name string
		args args
		want *Downloader
	}{
		{
			name: "positive testing",
			args: args{
				filepathProxy:      filepathProxy,
				gzipProxy:          gzipProxy,
				httpProxy:          httpProxy,
				ioProxy:            ioProxy,
				osProxy:            osProxy,
				utility:            util,
				wnJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			want: &Downloader{
				FilepathProxy: filepathProxy,
				GzipProxy:     gzipProxy,
				HttpProxy:     httpProxy,
				IoProxy:       ioProxy,
				OsProxy:       osProxy,
				Utility:       util,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.filepathProxy, tt.args.gzipProxy, tt.args.httpProxy, tt.args.ioProxy, tt.args.osProxy, tt.args.utility); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want=\n%v", got, tt.want)
			}
		})
	}
}

func TestDownloader_DownloadWNJpnDBFile(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	util := utility.New(
		fmtproxy.New(),
		osProxy,
		strconvproxy.New(),
	)
	dl := New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, WNJPN_DB_FILE_NAME)

	type fields struct {
		FilepathProxy      filepathproxy.FilePath
		GzipProxy          gzipproxy.Gzip
		HttpProxy          httpproxy.Http
		IoProxy            ioproxy.Io
		OsProxy            osproxy.Os
		Utility            utility.UtilityInterface
		WNJpnDBFileDirPath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    DownloadStatus
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (DB file exists)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			want:    DownloadedAlready,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (DB file does not exists)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			want:    DownloadedSuccessfully,
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
			name: "negative testing (Utility.CreateDirIfNotExits() failed)",
			fields: fields{
				FilepathProxy:      filepathproxy.New(),
				GzipProxy:          gzipproxy.New(),
				HttpProxy:          httpproxy.New(),
				IoProxy:            ioproxy.New(),
				OsProxy:            osproxy.New(),
				Utility:            nil,
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			want:    DownloadedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockUtil := mockutility.NewMockUtilityInterface(mockCtrl)
				mockUtil.EXPECT().CreateDirIfNotExist(gomock.Any()).Return(errors.New("Utility.CreateDirIfNotExist() failed"))
				tt.Utility = mockUtil
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
			d := New(
				tt.fields.FilepathProxy,
				tt.fields.GzipProxy,
				tt.fields.HttpProxy,
				tt.fields.IoProxy,
				tt.fields.OsProxy,
				tt.fields.Utility,
			)
			got, err := d.DownloadWNJpnDBFile(wnJpnDBFileDirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Downloader.DownloadWNJpnDBFile() : got = %v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestDownloader_downloadAndExtractDBFile(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, WNJPN_DB_FILE_NAME)

	type fields struct {
		FilepathProxy      filepathproxy.FilePath
		GzipProxy          gzipproxy.Gzip
		HttpProxy          httpproxy.Http
		IoProxy            ioproxy.Io
		OsProxy            osproxy.Os
		Utility            utility.UtilityInterface
		WNJpnDBFileDirPath string
	}
	type args struct {
		dbFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DownloadStatus
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				dbFilePath: wnJpnDBFilePath,
			},
			want:    DownloadedSuccessfully,
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
			name: "negative testing (Downloader.downloadGzipFile() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     nil,
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				dbFilePath: wnJpnDBFilePath,
			},
			want:    DownloadedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockHttpProxy := mockhttpproxy.NewMockHttp(mockCtrl)
				mockHttpProxy.EXPECT().Get(gomock.Any()).Return(nil, errors.New("HttpProxy.Get() failed"))
				tt.HttpProxy = mockHttpProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (Downloader.saveToTempFile() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       nil,
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				dbFilePath: wnJpnDBFilePath,
			},
			want:    DownloadedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockIoProxy := mockioproxy.NewMockIo(mockCtrl)
				mockIoProxy.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("IoProxy.Copy() failed"))
				tt.IoProxy = mockIoProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (Downloader.extractGzipFile() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     nil,
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				dbFilePath: wnJpnDBFilePath,
			},
			want:    DownloadedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockGzipProxy := mockgzipproxy.NewMockGzip(mockCtrl)
				mockGzipProxy.EXPECT().NewReader(gomock.Any()).Return(nil, errors.New("GzipProxy.NewReader() failed"))
				tt.GzipProxy = mockGzipProxy
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
			d := New(
				tt.fields.FilepathProxy,
				tt.fields.GzipProxy,
				tt.fields.HttpProxy,
				tt.fields.IoProxy,
				tt.fields.OsProxy,
				tt.fields.Utility,
			)
			got, err := d.downloadAndExtractDBFile(tt.args.dbFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.downloadAndExtractDBFile() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Downloader.downloadAndExtractDBFile() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestDownloader_downloadGzipFile(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}

	type fields struct {
		FilepathProxy      filepathproxy.FilePath
		GzipProxy          gzipproxy.Gzip
		HttpProxy          httpproxy.Http
		IoProxy            ioproxy.Io
		OsProxy            osproxy.Os
		Utility            utility.UtilityInterface
		WNJpnDBFileDirPath string
	}
	tests := []struct {
		name           string
		fields         fields
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "positive testing",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			wantStatusCode: 200,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New(
				tt.fields.FilepathProxy,
				tt.fields.GzipProxy,
				tt.fields.HttpProxy,
				tt.fields.IoProxy,
				tt.fields.OsProxy,
				tt.fields.Utility,
			)
			got, err := d.downloadGzipFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.downloadGzipFile() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got.FieldResponse.StatusCode != tt.wantStatusCode {
				t.Errorf("Downloader.downloadGzipFile() : got =\n%v, want =\n%v", got.FieldResponse.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestDownloader_saveToTempFile(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	util := utility.New(
		fmtproxy.New(),
		osProxy,
		strconvproxy.New(),
	)
	dl := New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	tempFilePath := filepathProxy.Join(osProxy.TempDir(), WNJPN_DB_ARCHIVE_FILE_NAME)
	resp, err := dl.downloadGzipFile()
	if err != nil {
		t.Errorf("Downloader.downloadGzipFile() : error =\n%v", err)
	}

	type fields struct {
		FilepathProxy      filepathproxy.FilePath
		GzipProxy          gzipproxy.Gzip
		HttpProxy          httpproxy.Http
		IoProxy            ioproxy.Io
		OsProxy            osproxy.Os
		Utility            utility.UtilityInterface
		WNJpnDBFileDirPath string
	}
	type args struct {
		body ioproxy.ReaderInstanceInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				body: resp.FieldResponse.Body,
			},
			want:    tempFilePath,
			wantErr: false,
			setup:   nil,
		},
		{
			name: "negative testing (OsProxy.Create() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				body: resp.FieldResponse.Body,
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().TempDir().Return(osProxy.TempDir())
				mockOsProxy.EXPECT().Create(gomock.Any()).Return(nil, errors.New("OsProxy.Create() failed"))
				tt.OsProxy = mockOsProxy
			},
		},
		{
			name: "negative testing (IoProxy.Copy() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       nil,
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				body: resp.FieldResponse.Body,
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockIoProxy := mockioproxy.NewMockIo(mockCtrl)
				mockIoProxy.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("IoProxy.Copy() failed"))
				tt.IoProxy = mockIoProxy
			},
		},
		{
			name: "negative testing (out.Seek() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				body: resp.FieldResponse.Body,
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().TempDir().Return(osProxy.TempDir())
				mockFileInstance := mockosproxy.NewMockFileInstanceInterface(mockCtrl)
				mockFileInstance.EXPECT().Seek(int64(0), ioproxy.SeekStart).Return(int64(0), errors.New("FileInstance.Seek() failed"))
				mockFileInstance.EXPECT().Close().Return(nil)
				mockOsProxy.EXPECT().Create(gomock.Any()).Return(mockFileInstance, nil)
				tt.OsProxy = mockOsProxy
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
				tt.fields.GzipProxy,
				tt.fields.HttpProxy,
				tt.fields.IoProxy,
				tt.fields.OsProxy,
				tt.fields.Utility,
			)
			got, err := d.saveToTempFile(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.saveToTempFile() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Downloader.saveToTempFile() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestDownloader_extractGzipFile(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, WNJPN_DB_FILE_NAME)
	util := utility.New(
		fmtproxy.New(),
		osProxy,
		strconvproxy.New(),
	)
	dl := New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	srcPath := filepathProxy.Join(osProxy.TempDir(), WNJPN_DB_ARCHIVE_FILE_NAME)
	destPath := filepathProxy.Join(wnJpnDBFileDirPath, WNJPN_DB_FILE_NAME)

	type fields struct {
		FilepathProxy      filepathproxy.FilePath
		GzipProxy          gzipproxy.Gzip
		HttpProxy          httpproxy.Http
		IoProxy            ioproxy.Io
		OsProxy            osproxy.Os
		Utility            utility.UtilityInterface
		WNJpnDBFileDirPath string
	}
	type args struct {
		srcPath  string
		destPath string
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
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				srcPath:  srcPath,
				destPath: destPath,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				resp, err := dl.downloadGzipFile()
				if err != nil {
					t.Errorf("Downloader.downloadGzipFile() : error =\n%v", err)
				}
				_, err = dl.saveToTempFile(resp.FieldResponse.Body)
				if err != nil {
					t.Errorf("Downloader.saveToTempFile() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (OsProxy.Open() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				srcPath:  srcPath,
				destPath: destPath,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				resp, err := dl.downloadGzipFile()
				if err != nil {
					t.Errorf("Downloader.downloadGzipFile() : error =\n%v", err)
				}
				_, err = dl.saveToTempFile(resp.FieldResponse.Body)
				if err != nil {
					t.Errorf("Downloader.saveToTempFile() : error =\n%v", err)
				}
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().Open(gomock.Any()).Return(nil, errors.New("OsProxy.Open() failed"))
				tt.OsProxy = mockOsProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (GzipProxy.NewReader() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     nil,
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				srcPath:  srcPath,
				destPath: destPath,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				resp, err := dl.downloadGzipFile()
				if err != nil {
					t.Errorf("Downloader.downloadGzipFile() : error =\n%v", err)
				}
				_, err = dl.saveToTempFile(resp.FieldResponse.Body)
				if err != nil {
					t.Errorf("Downloader.saveToTempFile() : error =\n%v", err)
				}
				mockGzipProxy := mockgzipproxy.NewMockGzip(mockCtrl)
				mockGzipProxy.EXPECT().NewReader(gomock.Any()).Return(nil, errors.New("GzipProxy.NewReader() failed"))
				tt.GzipProxy = mockGzipProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (OsProxy.Create() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       ioproxy.New(),
				OsProxy:       nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				srcPath:  srcPath,
				destPath: destPath,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				resp, err := dl.downloadGzipFile()
				if err != nil {
					t.Errorf("Downloader.downloadGzipFile() : error =\n%v", err)
				}
				_, err = dl.saveToTempFile(resp.FieldResponse.Body)
				if err != nil {
					t.Errorf("Downloader.saveToTempFile() : error =\n%v", err)
				}
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				fileInstance, err := osProxy.Open(srcPath)
				if err != nil {
					t.Errorf("OsProxy.Open() : error =\n%v", err)
				}
				mockOsProxy.EXPECT().Open(gomock.Any()).Return(fileInstance, nil)
				mockOsProxy.EXPECT().Create(gomock.Any()).Return(nil, errors.New("OsProxy.Create() failed"))
				tt.OsProxy = mockOsProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (IoProxy.Copy() failed)",
			fields: fields{
				FilepathProxy: filepathproxy.New(),
				GzipProxy:     gzipproxy.New(),
				HttpProxy:     httpproxy.New(),
				IoProxy:       nil,
				OsProxy:       osproxy.New(),
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				WNJpnDBFileDirPath: wnJpnDBFileDirPath,
			},
			args: args{
				srcPath:  srcPath,
				destPath: destPath,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				resp, err := dl.downloadGzipFile()
				if err != nil {
					t.Errorf("Downloader.downloadGzipFile() : error =\n%v", err)
				}
				_, err = dl.saveToTempFile(resp.FieldResponse.Body)
				if err != nil {
					t.Errorf("Downloader.saveToTempFile() : error =\n%v", err)
				}
				mockIoProxy := mockioproxy.NewMockIo(mockCtrl)
				mockIoProxy.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("IoProxy.Copy() failed"))
				tt.IoProxy = mockIoProxy
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
			d := New(
				tt.fields.FilepathProxy,
				tt.fields.GzipProxy,
				tt.fields.HttpProxy,
				tt.fields.IoProxy,
				tt.fields.OsProxy,
				tt.fields.Utility,
			)
			if err := d.extractGzipFile(tt.args.srcPath, tt.args.destPath); (err != nil) != tt.wantErr {
				t.Errorf("Downloader.extractGzipFile() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
