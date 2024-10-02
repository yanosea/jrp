package generator

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/database/wnjpn/model"
	"github.com/yanosea/jrp/app/database/wnjpn/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/downloader"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/gzip"
	"github.com/yanosea/jrp/app/proxy/http"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/rand"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/time"
	"github.com/yanosea/jrp/app/proxy/user"

	"github.com/yanosea/jrp/mock/app/database/wnjpn/repository"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	osProxy := osproxy.New()
	randProxy := randproxy.New()
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	wnJpnRepository := repository.New(
		sqlProxy,
	)

	type args struct {
		osProxy         osproxy.Os
		randProxy       randproxy.Rand
		sqlProxy        sqlproxy.Sql
		timeProxy       timeproxy.Time
		wnJpnRepository repository.WNJpnRepositoryInterface
	}
	tests := []struct {
		name string
		args args
		want *Generator
	}{
		{
			name: "positive testing",
			args: args{
				osProxy:         osProxy,
				randProxy:       randProxy,
				sqlProxy:        sqlProxy,
				timeProxy:       timeProxy,
				wnJpnRepository: wnJpnRepository,
			},
			want: &Generator{
				OsProxy:         osProxy,
				RandProxy:       randProxy,
				SqlProxy:        sqlProxy,
				TimeProxy:       timeProxy,
				WNJpnRepository: wnJpnRepository,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.osProxy, tt.args.randProxy, tt.args.sqlProxy, tt.args.timeProxy, tt.args.wnJpnRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want=\n%v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GenerateJrp(t *testing.T) {
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

	type fields struct {
		OsProxy         osproxy.Os
		RandProxy       randproxy.Rand
		SqlProxy        sqlproxy.Sql
		TimeProxy       timeproxy.Time
		WNJpnRepository repository.WNJpnRepositoryInterface
	}
	type args struct {
		wnJpnDBFilePath string
		num             int
		word            string
		mode            GenerateMode
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult GenerateResult
		wantCount  int
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller, tt *fields)
		cleanup    func()
	}{
		{
			name: "positive testing",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				num:             1,
				word:            "",
				mode:            WithNoPrefixOrSuffix,
			},
			wantResult: GeneratedSuccessfully,
			wantCount:  1,
			wantErr:    false,
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
			name: "positive testing (DB file not exists)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				num:             1,
				word:            "",
				mode:            WithNoPrefixOrSuffix,
			},
			wantResult: DBFileNotFound,
			wantCount:  0,
			wantErr:    false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: nil,
		},
		{
			name: "negative testing (Generator.getAllWords() failed)",
			fields: fields{
				OsProxy:         osproxy.New(),
				RandProxy:       randproxy.New(),
				SqlProxy:        sqlproxy.New(),
				TimeProxy:       timeproxy.New(),
				WNJpnRepository: nil,
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				num:             1,
				word:            "",
				mode:            WithNoPrefixOrSuffix,
			},
			wantResult: GeneratedFailed,
			wantCount:  0,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockWNJpnRepository := mockrepository.NewMockWNJpnRepositoryInterface(mockCtrl)
				mockWNJpnRepository.EXPECT().GetAllAVNWords(gomock.Any()).Return(nil, errors.New("WNJpnRepository.GetAllAVNWords() failed"))
				tt.WNJpnRepository = mockWNJpnRepository
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
			g := New(
				tt.fields.OsProxy,
				tt.fields.RandProxy,
				tt.fields.SqlProxy,
				tt.fields.TimeProxy,
				tt.fields.WNJpnRepository,
			)
			gotResult, jrps, err := g.GenerateJrp(tt.args.wnJpnDBFilePath, tt.args.num, tt.args.word, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.GenerateJrp() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Generator.GenerateJrp() : got =\n%v, want =\n%v", gotResult, tt.wantResult)
			}
			if len(jrps) != tt.wantCount {
				t.Errorf("Generator.GenerateJrp() : got(count) =\n%v, want =\n%v", len(jrps), tt.wantCount)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestGenerator_getAllWords(t *testing.T) {
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
	wnJpnRepository := repository.New(
		sqlproxy.New(),
	)
	if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
		t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
	}
	allAVNWords, err := wnJpnRepository.GetAllAVNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVNWords() : error =\n%v", err)
	}
	allAVWords, err := wnJpnRepository.GetAllAVWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVWords() : error =\n%v", err)
	}
	allNWords, err := wnJpnRepository.GetAllNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllNWords() : error =\n%v", err)
	}

	type fields struct {
		OsProxy         osproxy.Os
		RandProxy       randproxy.Rand
		SqlProxy        sqlproxy.Sql
		TimeProxy       timeproxy.Time
		WNJpnRepository repository.WNJpnRepositoryInterface
	}
	type args struct {
		wnJpnDBFilePath string
		mode            GenerateMode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Word
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (mode is NoPrefixOrSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				mode:            WithNoPrefixOrSuffix,
			},
			want:    allAVNWords,
			wantErr: false,
			setup: func() {
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
			name: "positive testing (mode is WithPrefix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				mode:            WithPrefix,
			},
			want:    allNWords,
			wantErr: false,
			setup: func() {
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
			name: "positive testing (mode is WithSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				mode:            WithSuffix,
			},
			want:    allAVWords,
			wantErr: false,
			setup: func() {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			g := New(
				tt.fields.OsProxy,
				tt.fields.RandProxy,
				tt.fields.SqlProxy,
				tt.fields.TimeProxy,
				tt.fields.WNJpnRepository,
			)
			got, err := g.getAllWords(tt.args.wnJpnDBFilePath, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.getAllWords() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generator.getAllWords() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
	if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
		t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
	}
}

func TestGenerator_getJrps(t *testing.T) {
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
	wnJpnRepository := repository.New(
		sqlproxy.New(),
	)
	if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
		t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
	}
	allAVWords, err := wnJpnRepository.GetAllAVWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVWords() : error =\n%v", err)
	}
	allNWords, err := wnJpnRepository.GetAllNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllNWords() : error =\n%v", err)
	}

	type fields struct {
		OsProxy         osproxy.Os
		RandProxy       randproxy.Rand
		SqlProxy        sqlproxy.Sql
		TimeProxy       timeproxy.Time
		WNJpnRepository repository.WNJpnRepositoryInterface
	}
	type args struct {
		num        int
		allAVWords []model.Word
		allNWords  []model.Word
		prefix     string
		suffix     string
		mode       GenerateMode
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantCount  int
		wantPrefix string
		wantSuffix string
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing (num is -1, mode is NoPrefixOrSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        0,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithNoPrefixOrSuffix,
			},
			wantCount:  0,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is -1, mode is WithPrefix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        -1,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithPrefix,
			},
			wantCount:  0,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is -1, mode is WithSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        -1,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithSuffix,
			},
			wantCount:  0,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 0, mode is NoPrefixOrSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        0,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithNoPrefixOrSuffix,
			},
			wantCount:  0,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 0, mode is WithPrefix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        0,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithPrefix,
			},
			wantCount:  0,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 0, mode is WithSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        0,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithSuffix,
			},
			wantCount:  0,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 1, mode is NoPrefixOrSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        1,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithNoPrefixOrSuffix,
			},
			wantCount:  1,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 1, mode is WithPrefix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        1,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "testPrefix",
				suffix:     "dummySuffix",
				mode:       WithPrefix,
			},
			wantCount:  1,
			wantPrefix: "testPrefix",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 1, mode is WithSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        1,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "testSuffix",
				mode:       WithSuffix,
			},
			wantCount:  1,
			wantPrefix: "",
			wantSuffix: "testSuffix",
			setup: func() {
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
			name: "positive testing (num is 2, mode is NoPrefixOrSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        2,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "dummySuffix",
				mode:       WithNoPrefixOrSuffix,
			},
			wantCount:  2,
			wantPrefix: "",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 2, mode is WithPrefix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        2,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "testPrefix",
				suffix:     "dummySuffix",
				mode:       WithPrefix,
			},
			wantCount:  2,
			wantPrefix: "testPrefix",
			wantSuffix: "",
			setup: func() {
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
			name: "positive testing (num is 2, mode is WithSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				num:        2,
				allAVWords: allAVWords,
				allNWords:  allNWords,
				prefix:     "dummyPrefix",
				suffix:     "testSuffix",
				mode:       WithSuffix,
			},
			wantCount:  2,
			wantPrefix: "",
			wantSuffix: "testSuffix",
			setup: func() {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			g := New(
				tt.fields.OsProxy,
				tt.fields.RandProxy,
				tt.fields.SqlProxy,
				tt.fields.TimeProxy,
				tt.fields.WNJpnRepository,
			)
			gots := g.getJrps(tt.args.num, tt.args.allAVWords, tt.args.allNWords, tt.args.prefix, tt.args.suffix, tt.args.mode)
			if len(gots) != tt.wantCount {
				t.Errorf("Generator.getJrps() : got(count) =\n%v, want =\n%v", len(gots), tt.wantCount)
			}
			for _, got := range gots {
				if got.Prefix.FieldNullString.String != tt.wantPrefix {
					t.Errorf("Generator.getJrps() : got(prefix) =\n%v, want =\n%v", gots[0].Prefix.FieldNullString.String, tt.wantPrefix)
				}
				if got.Suffix.FieldNullString.String != tt.wantSuffix {
					t.Errorf("Generator.getJrps() : got(suffix) =\n%v, want =\n%v", gots[0].Suffix.FieldNullString.String, tt.wantSuffix)
				}
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
		if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
			t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
		}
	}
}

func TestGenerator_getPrefixAndSuffix(t *testing.T) {
	type fields struct {
		OsProxy         osproxy.Os
		RandProxy       randproxy.Rand
		SqlProxy        sqlproxy.Sql
		TimeProxy       timeproxy.Time
		WNJpnRepository repository.WNJpnRepositoryInterface
	}
	type args struct {
		word string
		mode GenerateMode
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantPrefix string
		wantSuffix string
	}{
		{
			name: "positive testing (mode is NoPrefixOrSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				word: "dummy",
				mode: WithNoPrefixOrSuffix,
			},
			wantPrefix: "",
			wantSuffix: "",
		},
		{
			name: "positive testing (mode is WithPrefix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				word: "prefix",
				mode: WithPrefix,
			},
			wantPrefix: "prefix",
			wantSuffix: "",
		},
		{
			name: "positive testing (mode is WithSuffix)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				word: "suffix",
				mode: WithSuffix,
			},
			wantPrefix: "",
			wantSuffix: "suffix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(
				tt.fields.OsProxy,
				tt.fields.RandProxy,
				tt.fields.SqlProxy,
				tt.fields.TimeProxy,
				tt.fields.WNJpnRepository,
			)
			gotPrefix, gotSuffix := g.getPrefixAndSuffix(tt.args.word, tt.args.mode)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("Generator.getPrefixAndSuffix() : gotPrefix =\n%v, want =\n%v", gotPrefix, tt.wantPrefix)
			}
			if gotSuffix != tt.wantSuffix {
				t.Errorf("Generator.getPrefixAndSuffix() : gotSuffix =\n%v, want =\n%v", gotSuffix, tt.wantSuffix)
			}
		})
	}
}

func TestGenerator_separateWords(t *testing.T) {
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
	wnJpnRepository := repository.New(
		sqlproxy.New(),
	)
	if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
		t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
	}
	allAVNWords, err := wnJpnRepository.GetAllAVNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVNWords() : error =\n%v", err)
	}
	allAVWords, err := wnJpnRepository.GetAllAVWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllAVWords() : error =\n%v", err)
	}
	allNWords, err := wnJpnRepository.GetAllNWords(wnJpnDBFilePath)
	if err != nil {
		t.Errorf("WNJpnRepository.GetAllNWords() : error =\n%v", err)
	}
	type fields struct {
		OsProxy         osproxy.Os
		RandProxy       randproxy.Rand
		SqlProxy        sqlproxy.Sql
		TimeProxy       timeproxy.Time
		WNJpnRepository repository.WNJpnRepositoryInterface
	}
	type args struct {
		allWords []model.Word
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantAVWords []model.Word
		wantNWords  []model.Word
		setup       func()
		cleanup     func()
	}{
		{
			name: "positive testing (allWords is nil)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				allWords: nil,
			},
			wantAVWords: []model.Word{},
			wantNWords:  []model.Word{},
			setup: func() {
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
			name: "positive testing (allWords is empty)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				allWords: []model.Word{},
			},
			wantAVWords: []model.Word{},
			wantNWords:  []model.Word{},
			setup: func() {
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
			name: "positive testing (allWords is allAVNWords)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				allWords: allAVNWords,
			},
			wantAVWords: allAVWords,
			wantNWords:  allNWords,
			setup: func() {
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
			name: "positive testing (allWords is allAVWords)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				allWords: allAVWords,
			},
			wantAVWords: allAVWords,
			wantNWords:  []model.Word{},
			setup: func() {
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
			name: "positive testing (allWords is allNWords)",
			fields: fields{
				OsProxy:   osproxy.New(),
				RandProxy: randproxy.New(),
				SqlProxy:  sqlproxy.New(),
				TimeProxy: timeproxy.New(),
				WNJpnRepository: repository.New(
					sqlproxy.New(),
				),
			},
			args: args{
				allWords: allNWords,
			},
			wantAVWords: []model.Word{},
			wantNWords:  allNWords,
			setup: func() {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			g := New(
				tt.fields.OsProxy,
				tt.fields.RandProxy,
				tt.fields.SqlProxy,
				tt.fields.TimeProxy,
				tt.fields.WNJpnRepository,
			)
			avWords, nWords := g.separateWords(tt.args.allWords)
			if !reflect.DeepEqual(avWords, tt.wantAVWords) {
				t.Errorf("Generator.separateWords() : got =\n%v, want =\n%v", avWords, tt.wantAVWords)
			}
			if !reflect.DeepEqual(nWords, tt.wantNWords) {
				t.Errorf("Generator.separateWords() : got =\n%v, want =\n%v", nWords, tt.wantNWords)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
