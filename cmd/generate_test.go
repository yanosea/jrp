package cmd

import (
	"errors"
	"testing"

	"github.com/yanosea/jrp/app/database/jrp/model"
	jrprepository "github.com/yanosea/jrp/app/database/jrp/repository"
	wnjpnrepository "github.com/yanosea/jrp/app/database/wnjpn/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/downloader"
	"github.com/yanosea/jrp/app/library/generator"
	"github.com/yanosea/jrp/app/library/jrpwriter"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/gzip"
	"github.com/yanosea/jrp/app/proxy/http"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/rand"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/tablewriter"
	"github.com/yanosea/jrp/app/proxy/time"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"

	mockjrprepository "github.com/yanosea/jrp/mock/app/database/jrp/repository"
	"github.com/yanosea/jrp/mock/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/mock/app/library/generator"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNewGenerateCommand(t *testing.T) {
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
			got := NewGenerateCommand(tt.args.g)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewGenerateCommand().Execute() error =\n%v, wantError =\n%v", err, tt.wantError)
			}
		})
	}
}

func Test_generateOption_generateRunE(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	colorProxy := colorproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	wnJpnRepository := wnjpnrepository.New(
		sqlProxy,
	)
	gen := generator.New(
		osProxy,
		randproxy.New(),
		sqlProxy,
		timeproxy.New(),
		wnJpnRepository,
	)
	fmtProxy := fmtproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	strconvProxy := strconvproxy.New()
	jrpWriter := jrpwriter.New(
		strconvProxy,
		tablewriterproxy.New(),
	)
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	filepathProxy := filepathproxy.New()
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDBFileDirPathProviderFailsWNJpnDBFileDirPath := mockdbfiledirpathprovider.NewMockDBFileDirPathProvidable(mockCtrl)
	mockDBFileDirPathProviderFailsWNJpnDBFileDirPath.EXPECT().GetWNJpnDBFileDirPath().Return("", errors.New("DBFileDirPathProvider.GetWNJpnDBFileDirPath() failed"))
	mockDBFileDirPathProviderFailsGetJrpDBFileDirPath := mockdbfiledirpathprovider.NewMockDBFileDirPathProvidable(mockCtrl)
	mockDBFileDirPathProviderFailsGetJrpDBFileDirPath.EXPECT().GetWNJpnDBFileDirPath().Return("", nil)
	mockDBFileDirPathProviderFailsGetJrpDBFileDirPath.EXPECT().GetJrpDBFileDirPath().Return("", errors.New("DBFileDirPathProvider.GetJrpDBFileDirPath() failed"))

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
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing (wn jpn database file does not exist)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED) + "\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "positive testing (wn jpn database file exists)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: testutility.TEST_OUTPUT_ANY,
			wantStdErr: "",
			wantErr:    false,
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
			name: "positive testing (both prefix and suffix are specified)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "prefix",
						Suffix:                "suffix",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_USE_ONLY_ONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "positive testing (only prefix is specified)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "prefix",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: testutility.TEST_OUTPUT_ANY,
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "positive testing (only suffix is specified)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "suffix",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: testutility.TEST_OUTPUT_ANY,
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "positive testing (Args is nil)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED) + "\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "positive testing (Args is empty)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED) + "\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "negative testing (DBFileDirPathProvider.GetWNJpnDBFileDirPath() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: mockDBFileDirPathProviderFailsWNJpnDBFileDirPath,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "negative testing (DBFileDirPathProvider.GetJrpDBFileDirPath() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: mockDBFileDirPathProviderFailsGetJrpDBFileDirPath,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generateRunE(nil, nil)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if tt.wantStdOut != testutility.TEST_OUTPUT_ANY && stdout != tt.wantStdOut {
				t.Errorf("generateOption.generateRunE() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("generateOption.generateRunE() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_generateOption_generate(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	colorProxy := colorproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	wnJpnRepository := wnjpnrepository.New(
		sqlProxy,
	)
	gen := generator.New(
		osProxy,
		randproxy.New(),
		sqlProxy,
		timeproxy.New(),
		wnJpnRepository,
	)
	fmtProxy := fmtproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	strconvProxy := strconvproxy.New()
	jrpWriter := jrpwriter.New(
		strconvProxy,
		tablewriterproxy.New(),
	)
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	filepathProxy := filepathproxy.New()
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockGenerator := mockgenerator.NewMockGeneratable(mockCtrl)
	mockGenerator.EXPECT().GenerateJrp(wnJpnDBFilePath, 1, "", generator.WithNoPrefixOrSuffix).Return(generator.GeneratedFailed, nil, errors.New("Generator.GenerateJrp() failed"))
	mockJrpRepository := mockjrprepository.NewMockJrpRepositoryInterface(mockCtrl)
	mockJrpRepository.EXPECT().SaveHistory(gomock.Any(), gomock.Any()).Return(jrprepository.SavedFailed, errors.New("JrpRepository.SaveHistory() failed"))

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
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing (wn jpn database file does not exist)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generate(wnJpnDBFilePath, jrpDBFilePath, "", generator.WithNoPrefixOrSuffix)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED) + "\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
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
			name: "positive testing (wn jpn database file exists)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generate(wnJpnDBFilePath, jrpDBFilePath, "", generator.WithNoPrefixOrSuffix)
				},
				capturer: capturer,
			},
			wantStdOut: testutility.TEST_OUTPUT_ANY,
			wantStdErr: "",
			wantErr:    false,
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
			name: "negative testing (generateOption.generateGenerate() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             mockGenerator,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generate(wnJpnDBFilePath, jrpDBFilePath, "", generator.WithNoPrefixOrSuffix)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.GENERATE_MESSAGE_GENERATE_FAILURE) + "\n",
			wantErr:    false,
			setup: func() {
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
			name: "negative testing (generateOption.generateSave() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         mockJrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.generate(wnJpnDBFilePath, jrpDBFilePath, "", generator.WithNoPrefixOrSuffix)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.GENERATE_MESSAGE_SAVED_FAILURE) + "\n",
			wantErr:    false,
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
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if tt.wantStdOut != testutility.TEST_OUTPUT_ANY && stdout != tt.wantStdOut {
				t.Errorf("generateOption.generate() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("generateOption.generate() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_generateOption_generateGenerate(t *testing.T) {
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	wnJpnRepository := wnjpnrepository.New(
		sqlProxy,
	)
	gen := generator.New(
		osProxy,
		randproxy.New(),
		sqlProxy,
		timeproxy.New(),
		wnJpnRepository,
	)
	fmtProxy := fmtproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	strconvProxy := strconvproxy.New()
	jrpWriter := jrpwriter.New(
		strconvProxy,
		tablewriterproxy.New(),
	)
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	filepathProxy := filepathproxy.New()
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
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Args                  []string
		Number                int
		Prefix                string
		Suffix                string
		DryRun                bool
		Plain                 bool
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		Utility               utility.UtilityInterface
	}
	type args struct {
		wnJpnDBFilePath string
		word            string
		mode            generator.GenerateMode
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
		wantErr   bool
		setup     func(mockCtrl *gomock.Controller, tt *fields)
		cleanup   func()
	}{
		{
			name: "positive testing (wn jpn database file does not exist)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                false,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				word:            "dummy",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantCount: 0,
			wantErr:   false,
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
			name: "positive testing (wn jpn database file exists)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "prefix",
				Suffix:                "",
				DryRun:                false,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantCount: 1,
			wantErr:   false,
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
			name: "negative testing (Generator.GenerateJrp() failed)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "suffix",
				DryRun:                false,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             nil,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantCount: 0,
			wantErr:   true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockGenerator := mockgenerator.NewMockGeneratable(mockCtrl)
				mockGenerator.EXPECT().GenerateJrp(wnJpnDBFilePath, 1, "", generator.WithNoPrefixOrSuffix).Return(generator.GeneratedFailed, nil, errors.New("Generator.GenerateJrp() failed"))
				tt.Generator = mockGenerator
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
			o := &generateOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Args:                  tt.fields.Args,
				Number:                tt.fields.Number,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				DryRun:                tt.fields.DryRun,
				Plain:                 tt.fields.Plain,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				Utility:               tt.fields.Utility,
			}
			got, err := o.generateGenerate(tt.args.wnJpnDBFilePath, tt.args.word, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateOption.generateGenerate() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("generateOption.generateGenerate() : len(got) =\n%v, wantCount =\n%v", len(got), tt.wantCount)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_generateOption_writeGenerateGenerateResult(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	colorProxy := colorproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	wnJpnRepository := wnjpnrepository.New(
		sqlProxy,
	)
	gen := generator.New(
		osProxy,
		randproxy.New(),
		sqlProxy,
		timeproxy.New(),
		wnJpnRepository,
	)
	fmtProxy := fmtproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	strconvProxy := strconvproxy.New()
	jrpWriter := jrpwriter.New(
		strconvProxy,
		tablewriterproxy.New(),
	)
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)

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
			name: "positive testing (result is GeneratedSuccessfully)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateGenerateResult(generator.GeneratedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is DBFileNotFound)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateGenerateResult(generator.DBFileNotFound)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is GeneratedFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateGenerateResult(generator.GeneratedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.GENERATE_MESSAGE_GENERATE_FAILURE) + "\n",
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
			if tt.wantStdOut != testutility.TEST_OUTPUT_ANY && stdout != tt.wantStdOut {
				t.Errorf("generateOption.writeGenerateGenerateResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("generateOption.writeGenerateGenerateResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func Test_generateOption_generateSave(t *testing.T) {
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	wnJpnRepository := wnjpnrepository.New(
		sqlProxy,
	)
	gen := generator.New(
		osProxy,
		randproxy.New(),
		sqlProxy,
		timeproxy.New(),
		wnJpnRepository,
	)
	fmtProxy := fmtproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	strconvProxy := strconvproxy.New()
	jrpWriter := jrpwriter.New(
		strconvProxy,
		tablewriterproxy.New(),
	)
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	filepathProxy := filepathproxy.New()
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)
	timeProxy := timeproxy.New()
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Args                  []string
		Number                int
		Prefix                string
		Suffix                string
		DryRun                bool
		Plain                 bool
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		Utility               utility.UtilityInterface
	}
	type args struct {
		jrpDBFilePath string
		jrps          []model.Jrp
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantJrps []model.Jrp
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (not dry run, jrps are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                false,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps:          nil,
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not dry run, jrps are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                false,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps:          []model.Jrp{},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not dry run, jrps is one)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                false,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps: []model.Jrp{
					{
						Phrase:    "test1",
						Prefix:    sqlProxy.StringToNullString(""),
						Suffix:    sqlProxy.StringToNullString(""),
						CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
			},
			wantJrps: []model.Jrp{
				{
					ID:        1,
					Phrase:    "test1",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not dry run, jrps are two)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                false,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps: []model.Jrp{
					{
						Phrase:    "test1",
						Prefix:    sqlProxy.StringToNullString(""),
						Suffix:    sqlProxy.StringToNullString(""),
						CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
					{
						Phrase:    "test2",
						Prefix:    sqlProxy.StringToNullString(""),
						Suffix:    sqlProxy.StringToNullString(""),
						CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
			},
			wantJrps: []model.Jrp{
				{
					ID:        1,
					Phrase:    "test1",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:        2,
					Phrase:    "test2",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (dry run, jrps are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                true,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps:          nil,
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (dry run, jrps are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                true,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps:          []model.Jrp{},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (dry run, jrps is one)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                true,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps: []model.Jrp{
					{
						Phrase:    "test1",
						Prefix:    sqlProxy.StringToNullString(""),
						Suffix:    sqlProxy.StringToNullString(""),
						CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (dry run, jrps are two)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				Number:                1,
				Prefix:                "",
				Suffix:                "",
				DryRun:                true,
				Plain:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				Generator:             gen,
				JrpRepository:         jrpRepository,
				JrpWriter:             jrpWriter,
				WNJpnRepository:       wnJpnRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps: []model.Jrp{
					{
						Phrase:    "test1",
						Prefix:    sqlProxy.StringToNullString(""),
						Suffix:    sqlProxy.StringToNullString(""),
						CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
					{
						Phrase:    "test2",
						Prefix:    sqlProxy.StringToNullString(""),
						Suffix:    sqlProxy.StringToNullString(""),
						CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
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
			o := &generateOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Args:                  tt.fields.Args,
				Number:                tt.fields.Number,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				DryRun:                tt.fields.DryRun,
				Plain:                 tt.fields.Plain,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				Utility:               tt.fields.Utility,
			}
			if err := o.generateSave(tt.args.jrpDBFilePath, tt.args.jrps); (err != nil) != tt.wantErr {
				t.Errorf("generateOption.generateSave() : error = %v, wantErr %v", err, tt.wantErr)
			}
			savedJrps, err := jrpRepository.GetAllHistory(tt.args.jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(savedJrps, tt.wantJrps) {
				t.Errorf("generateOption.generateSave() : savedJrps =\n%v, want =\n%v", savedJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_generateOption_writeGenerateSaveResult(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	colorProxy := colorproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	wnJpnRepository := wnjpnrepository.New(
		sqlProxy,
	)
	gen := generator.New(
		osProxy,
		randproxy.New(),
		sqlProxy,
		timeproxy.New(),
		wnJpnRepository,
	)
	fmtProxy := fmtproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	strconvProxy := strconvproxy.New()
	jrpWriter := jrpwriter.New(
		strconvProxy,
		tablewriterproxy.New(),
	)
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)

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
			name: "positive testing (result is SavedSuccessfully)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateSaveResult(jrprepository.SavedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is SaveFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateSaveResult(jrprepository.SavedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.GENERATE_MESSAGE_SAVED_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (result is SavedNone)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateSaveResult(jrprepository.SavedNone)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_SAVED_NONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is SavedNotAll)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateSaveResult(jrprepository.SavedNotAll)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.GENERATE_MESSAGE_SAVED_NOT_ALL) + "\n",
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
			if tt.wantStdOut != testutility.TEST_OUTPUT_ANY && stdout != tt.wantStdOut {
				t.Errorf("generateOption.writeGenerateSaveResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("generateOption.writeGenerateSaveResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func Test_generateOption_writeGenerateResult(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	wnJpnRepository := wnjpnrepository.New(
		sqlProxy,
	)
	gen := generator.New(
		osProxy,
		randproxy.New(),
		sqlProxy,
		timeproxy.New(),
		wnJpnRepository,
	)
	fmtProxy := fmtproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	strconvProxy := strconvproxy.New()
	jrpWriter := jrpwriter.New(
		strconvProxy,
		tablewriterproxy.New(),
	)
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)
	timeProxy := timeproxy.New()

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
			name: "positive testing (jrps are nil, not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult(nil)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are empty, not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult([]model.Jrp{})
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps is one, not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult([]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString("prefix"),
							Suffix:    sqlProxy.StringToNullString("suffix"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					})
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\tPREFIX\tSUFFIX\tCREATED AT\ntest\tprefix\tsuffix\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 1\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps is two, not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult([]model.Jrp{
						{
							Phrase:    "test1",
							Prefix:    sqlProxy.StringToNullString("prefix1"),
							Suffix:    sqlProxy.StringToNullString("suffix1"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test2",
							Prefix:    sqlProxy.StringToNullString("prefix2"),
							Suffix:    sqlProxy.StringToNullString("suffix2"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					})
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\tPREFIX\tSUFFIX\tCREATED AT\ntest1\tprefix1\tsuffix1\t9999-12-31 00:00:00\ntest2\tprefix2\tsuffix2\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 2\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are nil, plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult(nil)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are empty, plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult([]model.Jrp{})
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps is one, plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult([]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString("prefix"),
							Suffix:    sqlProxy.StringToNullString("suffix"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					})
				},
				capturer: capturer,
			},
			wantStdOut: "test",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps is two, plain)",
			fields: fields{
				t: t,
				fnc: func() {
					generateOption := &generateOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Prefix:                "",
						Suffix:                "",
						DryRun:                false,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						Generator:             gen,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						WNJpnRepository:       wnJpnRepository,
						Utility:               util,
					}
					generateOption.writeGenerateResult([]model.Jrp{
						{
							Phrase:    "test1",
							Prefix:    sqlProxy.StringToNullString("prefix1"),
							Suffix:    sqlProxy.StringToNullString("suffix1"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test2",
							Prefix:    sqlProxy.StringToNullString("prefix2"),
							Suffix:    sqlProxy.StringToNullString("suffix2"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					})
				},
				capturer: capturer,
			},
			wantStdOut: "test1\ntest2",
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
			stdout = testutility.RemoveTabAndSpaceAndLf(stdout)
			stderr = testutility.RemoveTabAndSpaceAndLf(stderr)
			tt.wantStdOut = testutility.RemoveTabAndSpaceAndLf(tt.wantStdOut)
			tt.wantStdErr = testutility.RemoveTabAndSpaceAndLf(tt.wantStdErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if stdout != tt.wantStdOut {
				t.Errorf("generateOption.writeGenerateResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("generateOption.writeGenerateResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
