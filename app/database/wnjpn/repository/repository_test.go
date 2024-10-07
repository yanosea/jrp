package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/database/wnjpn/model"
	"github.com/yanosea/jrp/app/database/wnjpn/repository/query"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/downloader"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/gzip"
	"github.com/yanosea/jrp/app/proxy/http"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/user"

	"github.com/yanosea/jrp/mock/app/proxy/sql"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	sqlProxy := sqlproxy.New()

	type args struct {
		sqlProxy sqlproxy.Sql
	}
	tests := []struct {
		name string
		args args
		want *WNJpnRepository
	}{
		{
			name: "positive testing",
			args: args{
				sqlProxy: sqlProxy,
			},
			want: &WNJpnRepository{
				SqlProxy: sqlProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.sqlProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want=\n%v", got, tt.want)
			}
		})
	}
}

func TestWNJpnRepository_GetAllAVNWords(t *testing.T) {
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
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	wnJpnRepository := New(
		sqlproxy.New(),
	)
	if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
		t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
	}
	allAVNWords, err := wnJpnRepository.GetAllAVNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVNWords() : error =\n%v", err)
	}

	type fields struct {
		SqlProxy        sqlproxy.Sql
		WNJpnDBFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Word
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				SqlProxy: sqlproxy.New(),
			},
			want:    allAVNWords,
			wantErr: false,
			setup: func() {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.Remove(wnJpnDBFilePath); err != nil {
					t.Errorf("Os.Remove() : error =\n%v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			w := New(
				tt.fields.SqlProxy,
			)
			got, err := w.GetAllAVNWords(wnJpnDBFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("WNJpnRepository.GetAllAVNWords() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WNJpnRepository.GetAllAVNWords() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestWNJpnRepository_GetAllNWords(t *testing.T) {
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
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	wnJpnRepository := New(
		sqlproxy.New(),
	)
	if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
		t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
	}
	allNWords, err := wnJpnRepository.GetAllNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllNWords() : error =\n%v", err)
	}

	type fields struct {
		SqlProxy sqlproxy.Sql
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Word
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				SqlProxy: sqlproxy.New(),
			},
			want:    allNWords,
			wantErr: false,
			setup: func() {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.Remove(wnJpnDBFilePath); err != nil {
					t.Errorf("Os.Remove() : error =\n%v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			w := New(
				tt.fields.SqlProxy,
			)
			got, err := w.GetAllNWords(wnJpnDBFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("WNJpnRepository.GetAllNWords() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WNJpnRepository.GetAllNWords() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestWNJpnRepository_GetAllAVWords(t *testing.T) {
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
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	wnJpnRepository := New(
		sqlproxy.New(),
	)
	if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
		t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
	}
	allAVWords, err := wnJpnRepository.GetAllAVWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVWords() : error =\n%v", err)
	}

	type fields struct {
		SqlProxy        sqlproxy.Sql
		WNJpnDBFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Word
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				SqlProxy:        sqlproxy.New(),
				WNJpnDBFilePath: wnJpnDBFilePath,
			},
			want:    allAVWords,
			wantErr: false,
			setup: func() {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.Remove(wnJpnDBFilePath); err != nil {
					t.Errorf("Os.Remove() : error =\n%v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			w := New(
				tt.fields.SqlProxy,
			)
			got, err := w.GetAllAVWords(wnJpnDBFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("WNJpnRepository.GetAllAVWords() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WNJpnRepository.GetAllAVWords() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestWNJpnRepository_getWords(t *testing.T) {
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
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	wnJpnRepository := New(
		sqlproxy.New(),
	)
	if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
		t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
	}
	allAVNWords, err := wnJpnRepository.GetAllAVNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVNWords() : error =\n%v", err)
	}

	type fields struct {
		SqlProxy        sqlproxy.Sql
		WNJpnDBFilePath string
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Word
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				SqlProxy:        sqlproxy.New(),
				WNJpnDBFilePath: wnJpnDBFilePath,
			},
			args: args{
				query: query.GetAllJapaneseAVNWords,
			},
			want:    allAVNWords,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.Remove(wnJpnDBFilePath); err != nil {
					t.Errorf("Os.Remove() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				SqlProxy:        nil,
				WNJpnDBFilePath: wnJpnDBFilePath,
			},
			args: args{
				query: query.GetAllJapaneseAVNWords,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.Remove(wnJpnDBFilePath); err != nil {
					t.Errorf("Os.Remove() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Query() failed)",
			fields: fields{
				SqlProxy:        nil,
				WNJpnDBFilePath: wnJpnDBFilePath,
			},
			args: args{
				query: query.GetAllJapaneseAVNWords,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Query(gomock.Any()).Return(nil, errors.New("DBInstance.Query() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.Remove(wnJpnDBFilePath); err != nil {
					t.Errorf("Os.Remove() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				SqlProxy:        nil,
				WNJpnDBFilePath: wnJpnDBFilePath,
			},
			args: args{
				query: query.GetAllJapaneseAVNWords,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Query(gomock.Any()).Return(mockRowsInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.Remove(wnJpnDBFilePath); err != nil {
					t.Errorf("Os.Remove() : error =\n%v", err)
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
			w := New(
				tt.fields.SqlProxy,
			)
			got, err := w.getWords(wnJpnDBFilePath, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("WNJpnRepository.getWords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WNJpnRepository.getWords() = %v, want %v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
