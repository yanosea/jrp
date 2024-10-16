package cmd

import (
	"errors"
	"testing"

	"github.com/eiannone/keyboard"

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
	"github.com/yanosea/jrp/app/proxy/keyboard"
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
	"github.com/yanosea/jrp/mock/app/proxy/keyboard"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNewInteractiveCommand(t *testing.T) {
	osProxy := osproxy.New()
	filepathProxy := filepathproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
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
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)

	type args struct {
		g             *GlobalOption
		keyboardProxy keyboardproxy.Keyboard
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		setup     func(*gomock.Controller, *args)
		cleanup   func()
	}{
		{
			name: "positive testing",
			args: args{
				g: NewGlobalOption(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				keyboardProxy: nil,
			},
			wantError: false,
			setup: func(mockCtrl *gomock.Controller, args *args) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := ","
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(30).Return(r, keyboard.KeyEnter, nil)
				args.keyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
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
				tt.setup(mockCtrl, &tt.args)
			}
			got := NewInteractiveCommand(tt.args.g, tt.args.keyboardProxy)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewInteractiveCommand().Execute() : error =\n%v", err)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_interactiveOption_interactiveRunE(t *testing.T) {
	osProxy := osproxy.New()
	filepathProxy := filepathproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	fmtProxy := fmtproxy.New()
	strconvProxy := strconvproxy.New()
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Prefix                string
		Suffix                string
		Plain                 bool
		Timeout               int
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		KeyboardProxy         keyboardproxy.Keyboard
		Utility               utility.UtilityInterface
	}
	tests := []struct {
		name         string
		fields       fields
		wantJrpCount int
		wantErr      bool
		setup        func(*gomock.Controller, *fields)
		cleanup      func()
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantJrpCount: 1,
			wantErr:      false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "i"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (prefix is set)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "prefix",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantJrpCount: 1,
			wantErr:      false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "i"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (suffix is set)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "suffix",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantJrpCount: 1,
			wantErr:      false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "i"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (both prefix and suffix are set)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "prefix",
				Suffix:  "suffix",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantJrpCount: 0,
			wantErr:      false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (DBFileDirPathProvider.GetWNJpnDBFileDirPath() failed)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Prefix:                "",
				Suffix:                "",
				Plain:                 false,
				Timeout:               1,
				DBFileDirPathProvider: nil,
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantJrpCount: 0,
			wantErr:      true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockDBFileDirPathProvider := mockdbfiledirpathprovider.NewMockDBFileDirPathProvidable(mockCtrl)
				mockDBFileDirPathProvider.EXPECT().GetWNJpnDBFileDirPath().Return("", errors.New("DBFileDirPathProvider.GetWNJpnDBFileDirPath() failed"))
				fields.DBFileDirPathProvider = mockDBFileDirPathProvider
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (DBFileDirPathProvider.GetJrpDBFileDirPath() failed)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Prefix:                "",
				Suffix:                "",
				Plain:                 false,
				Timeout:               1,
				DBFileDirPathProvider: nil,
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantJrpCount: 0,
			wantErr:      true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockDBFileDirPathProvider := mockdbfiledirpathprovider.NewMockDBFileDirPathProvidable(mockCtrl)
				mockDBFileDirPathProvider.EXPECT().GetWNJpnDBFileDirPath().Return("", nil)
				mockDBFileDirPathProvider.EXPECT().GetJrpDBFileDirPath().Return("", errors.New("DBFileDirPathProvider.GetJrpDBFileDirPath() failed"))
				fields.DBFileDirPathProvider = mockDBFileDirPathProvider
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
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
			o := &interactiveOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				Plain:                 tt.fields.Plain,
				Timeout:               tt.fields.Timeout,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				KeyboardProxy:         tt.fields.KeyboardProxy,
				Utility:               tt.fields.Utility,
			}
			if err := o.interactiveRunE(nil, nil); (err != nil) != tt.wantErr {
				t.Errorf("interactiveOption.interactiveRunE() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			jrps, err := jrpRepository.GetAllHistory(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if len(jrps) != tt.wantJrpCount {
				t.Errorf("JrpRepository.GetAllHistory() : len(got) =\n%v, wantJrpCount =\n%v", len(jrps), tt.wantJrpCount)
			}
			if tt.wantJrpCount != 0 && tt.fields.Prefix != "" {
				for _, jrp := range jrps {
					if jrp.Prefix.FieldNullString.String != tt.fields.Prefix {
						t.Errorf("JrpRepository.GetAllHistory() : got.Prefix =\n%v, tt.fields.Prefix =\n%v", jrp.Prefix.FieldNullString.String, tt.fields.Prefix)
					}
				}
			}
			if tt.wantJrpCount != 0 && tt.fields.Suffix != "" {
				for _, jrp := range jrps {
					if jrp.Suffix.FieldNullString.String != tt.fields.Suffix {
						t.Errorf("JrpRepository.GetAllHistory() : got.Suffix =\n%v, tt.fields.Suffix =\n%v", jrp.Prefix.FieldNullString.String, tt.fields.Suffix)
					}
				}
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_interactiveOption_interactive(t *testing.T) {
	osProxy := osproxy.New()
	filepathProxy := filepathproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	fmtProxy := fmtproxy.New()
	strconvProxy := strconvproxy.New()
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvProxy,
	)
	dl := downloader.New(
		filepathProxy,
		gzipproxy.New(),
		httpproxy.New(),
		ioproxy.New(),
		osProxy,
		util,
	)
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Prefix                string
		Suffix                string
		Plain                 bool
		Timeout               int
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		KeyboardProxy         keyboardproxy.Keyboard
		Utility               utility.UtilityInterface
	}
	type args struct {
		wnJpnDBFilePath string
		jrpDBFilePath   string
		word            string
		mode            generator.GenerateMode
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantJrpCount          int
		wantFavoritedJrpCount int
		wantErr               bool
		setup                 func(*gomock.Controller, *fields)
		cleanup               func()
	}{
		{
			name: "positive testing (there is no word net japan db file)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing ('u' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          1,
			wantFavoritedJrpCount: 1,
			wantErr:               false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "u"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				ss := ","
				rr := rune(ss[0])
				mockKeyboardProxy.EXPECT().Close()
				mockKeyboardProxy.EXPECT().GetKey(1).Return(rr, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing ('i' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          1,
			wantFavoritedJrpCount: 1,
			wantErr:               false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "i"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing ('j' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          1,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "j"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				ss := ","
				rr := rune(ss[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(rr, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing ('k' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          1,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "k"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing ('m' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "m"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				ss := ","
				rr := rune(ss[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(rr, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (',' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := ","
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (InteractiveOption.getInteractiveInteractiveAnswer() failed)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(errors.New("KeyboardProxy.Open() failed"))
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (InteractiveOption.interactiveSave() failed)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: nil,
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "k"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockJrpRepository := mockjrprepository.NewMockJrpRepositoryInterface(mockCtrl)
				mockJrpRepository.EXPECT().SaveHistory(gomock.Any(), gomock.Any()).Return(jrprepository.SavedFailed, errors.New("JrpRepository.SaveHistory() failed"))
				fields.JrpRepository = mockJrpRepository
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (InteractiveOption.interactiveFavorite() failed)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: nil,
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				jrpDBFilePath:   jrpDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "i"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockJrpRepository := mockjrprepository.NewMockJrpRepositoryInterface(mockCtrl)
				mockJrpRepository.EXPECT().SaveHistory(gomock.Any(), gomock.Any()).Return(jrprepository.SavedSuccessfully, nil)
				mockJrpRepository.EXPECT().AddFavoriteByIDs(gomock.Any(), []int{0}).Return(jrprepository.AddedFailed, errors.New("JrpRepository.AddFavoriteByIDs() failed"))
				fields.JrpRepository = mockJrpRepository
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
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
			o := &interactiveOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				Plain:                 tt.fields.Plain,
				Timeout:               tt.fields.Timeout,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				KeyboardProxy:         tt.fields.KeyboardProxy,
				Utility:               tt.fields.Utility,
			}
			if err := o.interactive(
				tt.args.wnJpnDBFilePath,
				tt.args.jrpDBFilePath,
				tt.args.word,
				tt.args.mode,
			); (err != nil) != tt.wantErr {
				t.Errorf("interactiveOption.interactive() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			jrps, err := jrpRepository.GetAllHistory(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if len(jrps) != tt.wantJrpCount {
				t.Errorf("JrpRepository.GetAllHistory() : len(got) =\n%v, wantJrpCount =\n%v", len(jrps), tt.wantJrpCount)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if len(favoritedJrps) != tt.wantFavoritedJrpCount {
				t.Errorf("JrpRepository.GetAllFavorite() : len(got) =\n%v, wantFavoritedJrpCount =\n%v", len(favoritedJrps), tt.wantFavoritedJrpCount)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_interactiveOption_interactiveGenerate(t *testing.T) {
	osProxy := osproxy.New()
	filepathProxy := filepathproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
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
	wnJpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWNJpnDBFileDirPath() : error =\n%v", err)
	}
	wnJpnDBFilePath := filepathProxy.Join(wnJpnDBFileDirPath, downloader.WNJPN_DB_FILE_NAME)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Prefix                string
		Suffix                string
		Plain                 bool
		Timeout               int
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		KeyboardProxy         keyboardproxy.Keyboard
		Utility               utility.UtilityInterface
	}
	type args struct {
		wnJpnDBFilePath string
		word            string
		mode            generator.GenerateMode
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         generator.GenerateResult
		wantJrpCount int
		wantErr      bool
		setup        func()
		cleanup      func()
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				wnJpnDBFilePath: wnJpnDBFilePath,
				word:            "",
				mode:            generator.WithNoPrefixOrSuffix,
			},
			want:         generator.GeneratedSuccessfully,
			wantJrpCount: 1,
			wantErr:      false,
			setup: func() {
				if _, err := dl.DownloadWNJpnDBFile(wnJpnDBFileDirPath); err != nil {
					t.Errorf("Downloader.DownloadWNJpnDBFile() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(wnJpnDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
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
			o := &interactiveOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				Plain:                 tt.fields.Plain,
				Timeout:               tt.fields.Timeout,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				KeyboardProxy:         tt.fields.KeyboardProxy,
				Utility:               tt.fields.Utility,
			}
			gotJrps, gotResult, err := o.interactiveGenerate(tt.args.wnJpnDBFilePath, tt.args.word, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("interactiveOption.interactiveGenerate() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if len(gotJrps) != tt.wantJrpCount {
				t.Errorf("interactiveOption.interactiveGenerate() : len(gotJrps) =\n%v, wantJrpCount =\n%v", len(gotJrps), tt.wantJrpCount)
			}
			if gotResult != tt.want {
				t.Errorf("interactiveOption.interactiveGenerate() : gotResult =\n%v, want =\n%v", gotResult, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_interactiveOption_writeInteractiveGenerateResult(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
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
			name: "positive testing (arg is GeneratedSuccessfully)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveGenerateResult(generator.GeneratedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is GeneratedFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveGenerateResult(generator.GeneratedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.INTERACTIVE_MESSAGE_GENERATE_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is DBFileNotFound)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveGenerateResult(generator.DBFileNotFound)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED) + "\n",
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
				t.Errorf("interactiveOption.writeInteractiveGenerateResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("interactiveOption.writeInteractiveGenerateResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func Test_interactiveOption_writeInteractiveGeneratedJrp(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
	sqlProxy := sqlproxy.New()
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
			name: "positive testing (jrp is nil)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveGeneratedJrp(nil)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveGeneratedJrp(
						[]*model.Jrp{
							{
								Phrase:      "test",
								Prefix:      sqlProxy.StringToNullString(""),
								Suffix:      sqlProxy.StringToNullString(""),
								IsFavorited: 0,
								CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
								UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							},
						},
					)
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\tPREFIX\tSUFFIX\tCREATED AT\ntest\t\t\t\t9999-12-31 00:00:00\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (plain)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   true,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveGeneratedJrp(
						[]*model.Jrp{
							{
								Phrase:      "test",
								Prefix:      sqlProxy.StringToNullString(""),
								Suffix:      sqlProxy.StringToNullString(""),
								IsFavorited: 0,
								CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
								UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							},
						},
					)
				},
				capturer: capturer,
			},
			wantStdOut: "test\n\n",
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
			if tt.wantStdOut != testutility.TEST_OUTPUT_ANY && stdout != tt.wantStdOut {
				t.Errorf("interactiveOption.writeInteractiveGeneratedJrp() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("interactiveOption.writeInteractiveGeneratedJrp() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func Test_interactiveOption_interactiveSave(t *testing.T) {
	osProxy := osproxy.New()
	filepathProxy := filepathproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := jrprepository.New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Prefix                string
		Suffix                string
		Plain                 bool
		Timeout               int
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		KeyboardProxy         keyboardproxy.Keyboard
		Utility               utility.UtilityInterface
	}
	type args struct {
		jrpDBFilePath     string
		jrps              []*model.Jrp
		interactiveAnswer constant.InteractiveAnswer
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantJrpCount int
		wantErr      bool
		setup        func()
		cleanup      func()
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps: []*model.Jrp{
					{
						ID:          1,
						Phrase:      "test",
						Prefix:      sqlProxy.StringToNullString(""),
						Suffix:      sqlProxy.StringToNullString(""),
						IsFavorited: 0,
						CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
				interactiveAnswer: constant.InteractiveAnswerSaveAndExit,
			},
			wantJrpCount: 1,
			wantErr:      false,
			setup: func() {
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
				tt.setup()
			}
			o := &interactiveOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				Plain:                 tt.fields.Plain,
				Timeout:               tt.fields.Timeout,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				KeyboardProxy:         tt.fields.KeyboardProxy,
				Utility:               tt.fields.Utility,
			}
			err := o.interactiveSave(tt.args.jrpDBFilePath, tt.args.jrps, tt.args.interactiveAnswer)
			if (err != nil) != tt.wantErr {
				t.Errorf("interactiveOption.interactiveGenerate() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			savedJrps, err := jrpRepository.GetAllHistory(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if len(savedJrps) != tt.wantJrpCount {
				t.Errorf("interactiveOption.interactiveGenerate() : len(savedJrps) =\n%v, wantJrpCount =\n%v", len(savedJrps), tt.wantJrpCount)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_interactiveOption_writeInteractiveSaveResult(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
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
			name: "positive testing (arg is SavedSuccessfully, not favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedSuccessfully, constant.InteractiveAnswerSaveAndExit)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.INTERACTIVE_MESSAGE_SAVED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedSuccessfully, favorited, continue)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedSuccessfully, constant.InteractiveAnswerSaveAndFavoriteAndContinue)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedSuccessfully, favorited, exit)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedSuccessfully, constant.InteractiveAnswerSaveAndFavoriteAndExit)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedFailed, not favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedFailed, constant.InteractiveAnswerSaveAndExit)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.INTERACTIVE_MESSAGE_SAVED_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedFailed, favorited, continue)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedFailed, constant.InteractiveAnswerSaveAndContinue)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.INTERACTIVE_MESSAGE_SAVED_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedFailed, favorited, exit)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedFailed, constant.InteractiveAnswerSaveAndContinue)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.INTERACTIVE_MESSAGE_SAVED_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedNone, not favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedNone, constant.InteractiveAnswerSaveAndExit)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedNone, favorited, continue)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedNone, constant.InteractiveAnswerSaveAndContinue)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedNone, favorited, exit)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedNone, constant.InteractiveAnswerSaveAndContinue)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedNotAll, not favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedNotAll, constant.InteractiveAnswerSaveAndExit)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedNotAll, favorited, continue)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedNotAll, constant.InteractiveAnswerSaveAndContinue)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is SavedNotAll, favorited, exit)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveSaveResult(jrprepository.SavedNotAll, constant.InteractiveAnswerSaveAndContinue)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NOT_ALL) + "\n",
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
				t.Errorf("interactiveOption.writeInteractiveSavedResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("interactiveOption.writeInteractiveSavedResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func Test_interactiveOption_interactiveFavorite(t *testing.T) {
	osProxy := osproxy.New()
	filepathProxy := filepathproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := jrprepository.New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlProxy,
		stringsproxy.New(),
	)
	testJrps := []*model.Jrp{
		{
			ID:          1,
			Phrase:      "test",
			Prefix:      sqlProxy.StringToNullString(""),
			Suffix:      sqlProxy.StringToNullString(""),
			IsFavorited: 0,
			CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
			UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
		},
	}

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Prefix                string
		Suffix                string
		Plain                 bool
		Timeout               int
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		KeyboardProxy         keyboardproxy.Keyboard
		Utility               utility.UtilityInterface
	}
	type args struct {
		jrpDBFilePath string
		jrps          []*model.Jrp
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantFavoritedJrpCount int
		wantErr               bool
		setup                 func()
		cleanup               func()
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps:          testJrps,
			},
			wantFavoritedJrpCount: 1,
			wantErr:               false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(jrpDBFilePath, testJrps); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
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
				tt.setup()
			}
			o := &interactiveOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				Plain:                 tt.fields.Plain,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				KeyboardProxy:         tt.fields.KeyboardProxy,
				Utility:               tt.fields.Utility,
			}
			if err := o.interactiveFavorite(tt.args.jrpDBFilePath, tt.args.jrps); (err != nil) != tt.wantErr {
				t.Errorf("interactiveOption.interactiveFavorite() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if len(favoritedJrps) != tt.wantFavoritedJrpCount {
				t.Errorf("interactiveOption.interactiveFavorite() : len(favoritedJrps) =\n%v, wantFavoritedJrpCount =\n%v", len(favoritedJrps), tt.wantFavoritedJrpCount)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_interactiveOption_writeInteractiveFavoriteResult(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
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
			name: "positive testing (arg is AddedSuccessfully)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveFavoriteResult(jrprepository.AddedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.INTERACTIVE_MESSAGE_FAVORITED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is AddedFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveFavoriteResult(jrprepository.AddedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.INTERACTIVE_MESSAGE_FAVORITED_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is AddedNone)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveFavoriteResult(jrprepository.AddedNone)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_FAVORITED_NONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is AddedNotAll)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writeInteractiveFavoriteResult(jrprepository.AddedNotAll)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_FAVORITED_NOT_ALL) + "\n",
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
				t.Errorf("interactiveOption.writeInteractiveFavoriteResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("interactiveOption.writeInteractiveFavoriteResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func Test_interactiveOption_getInteractiveInteractiveAnswer(t *testing.T) {
	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Prefix                string
		Suffix                string
		Plain                 bool
		Timeout               int
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		Generator             generator.Generatable
		JrpRepository         jrprepository.JrpRepositoryInterface
		JrpWriter             jrpwriter.JrpWritable
		WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
		KeyboardProxy         keyboardproxy.Keyboard
		Utility               utility.UtilityInterface
	}
	type args struct {
		timeoutSec int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    constant.InteractiveAnswer
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, fields *fields)
	}{
		{
			name: "positive testing ('u' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSaveAndFavoriteAndContinue,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "u"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "positive testing ('i' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSaveAndFavoriteAndExit,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "i"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "positive testing ('j' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSaveAndContinue,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "j"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "positive testing ('k' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSaveAndExit,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "k"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "positive testing ('m' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSkipAndContinue,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := "m"
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "positive testing (',' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSkipAndExit,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := ","
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "positive testing (',' was inputted)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSkipAndExit,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				s := ","
				r := rune(s[0])
				mockKeyboardProxy.EXPECT().GetKey(1).Return(r, keyboard.KeyEnter, nil)
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "negative testing (KeyboardProxy.Open() failed)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSkipAndExit,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(errors.New("KeyboardProxy.Open() failed"))
				fields.KeyboardProxy = mockKeyboardProxy
			},
		},
		{
			name: "negative testing (KeyboardProxy.GetKey() failed)",
			fields: fields{
				Out:     osproxy.Stdout,
				ErrOut:  osproxy.Stderr,
				Prefix:  "",
				Suffix:  "",
				Plain:   false,
				Timeout: 1,
				DBFileDirPathProvider: dbfiledirpathprovider.New(
					filepathproxy.New(),
					osproxy.New(),
					userproxy.New(),
				),
				Generator: generator.New(
					osproxy.New(),
					randproxy.New(),
					sqlproxy.New(),
					timeproxy.New(),
					wnjpnrepository.New(
						sqlproxy.New(),
					),
				),
				JrpRepository: jrprepository.New(
					fmtproxy.New(),
					sortproxy.New(),
					sqlproxy.New(),
					stringsproxy.New(),
				),
				JrpWriter: jrpwriter.New(
					strconvproxy.New(),
					tablewriterproxy.New(),
				),
				WNJpnRepository: wnjpnrepository.New(
					sqlproxy.New(),
				),
				KeyboardProxy: nil,
				Utility: utility.New(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			args: args{
				timeoutSec: 1,
			},
			want:    constant.InteractiveAnswerSkipAndExit,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				mockKeyboardProxy := mockkeyboardproxy.NewMockKeyboard(mockCtrl)
				mockKeyboardProxy.EXPECT().Open().Return(nil)
				mockKeyboardProxy.EXPECT().Close()
				mockKeyboardProxy.EXPECT().GetKey(1).Return(rune(0), keyboard.KeyEnter, errors.New("KeyboardProxy.Getkey() failed"))
				fields.KeyboardProxy = mockKeyboardProxy
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
			o := &interactiveOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Prefix:                tt.fields.Prefix,
				Suffix:                tt.fields.Suffix,
				Plain:                 tt.fields.Plain,
				Timeout:               tt.fields.Timeout,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				Generator:             tt.fields.Generator,
				JrpRepository:         tt.fields.JrpRepository,
				JrpWriter:             tt.fields.JrpWriter,
				WNJpnRepository:       tt.fields.WNJpnRepository,
				KeyboardProxy:         tt.fields.KeyboardProxy,
				Utility:               tt.fields.Utility,
			}
			got, err := o.getInteractiveInteractiveAnswer(tt.args.timeoutSec)
			if (err != nil) != tt.wantErr {
				t.Errorf("interactiveOption.getInteractiveInteractiveAnswer() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("interactiveOption.getInteractiveInteractiveAnswer() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func Test_interactiveOption_writePhase(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
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
			name: "positive testing (arg is -1)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writePhase(-1)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.BlueString(constant.INTERACTIVE_MESSAGE_PHASE+"-1") + "\n\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is 0)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writePhase(0)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.BlueString(constant.INTERACTIVE_MESSAGE_PHASE+"0") + "\n\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (arg is 1)",
			fields: fields{
				t: t,
				fnc: func() {
					interactiveOption := &interactiveOption{
						Out:     osproxy.Stdout,
						ErrOut:  osproxy.Stderr,
						Prefix:  "",
						Suffix:  "",
						Plain:   false,
						Timeout: 1,
						DBFileDirPathProvider: dbfiledirpathprovider.New(
							filepathproxy.New(),
							osproxy.New(),
							userproxy.New(),
						),
						Generator: generator.New(
							osproxy.New(),
							randproxy.New(),
							sqlproxy.New(),
							timeproxy.New(),
							wnjpnrepository.New(
								sqlproxy.New(),
							),
						),
						JrpRepository: jrprepository.New(
							fmtproxy.New(),
							sortproxy.New(),
							sqlproxy.New(),
							stringsproxy.New(),
						),
						JrpWriter: jrpwriter.New(
							strconvproxy.New(),
							tablewriterproxy.New(),
						),
						WNJpnRepository: wnjpnrepository.New(
							sqlproxy.New(),
						),
						KeyboardProxy: nil,
						Utility: utility.New(
							fmtproxy.New(),
							osproxy.New(),
							strconvproxy.New(),
						),
					}
					interactiveOption.writePhase(1)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.BlueString(constant.INTERACTIVE_MESSAGE_PHASE+"1") + "\n\n",
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
				t.Errorf("interactiveOption.writePhase() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("interactiveOption.writePhase() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
