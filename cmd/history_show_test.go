package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/database/jrp/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/jrpwriter"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/tablewriter"
	"github.com/yanosea/jrp/app/proxy/time"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"

	"github.com/yanosea/jrp/mock/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/mock/app/library/utility"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNewHistoryShowCommand(t *testing.T) {
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
			got := NewHistoryShowCommand(tt.args.g)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewHistoryShowCommand().Execute() : error =\n%v", err)
			}
		})
	}
}

func Test_historyShowOption_historyShowRunE(t *testing.T) {
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
	fmtProxy := fmtproxy.New()
	jrpRepository := repository.New(
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
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME)
	timeProxy := timeproxy.New()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDBFileDirPathProvider := mockdbfiledirpathprovider.NewMockDBFileDirPathProvidable(mockCtrl)
	mockDBFileDirPathProvider.EXPECT().GetJrpDBFileDirPath().Return("", errors.New("DBFileDirPathProvider.GetJrpDBFileDirPath() failed"))
	mockUtility := mockutility.NewMockUtilityInterface(mockCtrl)
	mockUtility.EXPECT().CreateDirIfNotExist(jrpDBFileDirPath).Return(errors.New("Utility.CreateDirIfNotExist() failed"))

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
			name: "positive testing (no jrps in the database file)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShowRunE(nil, nil); err != nil {
						t.Errorf("historyShowOption.historyShowRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_MESSAGE_NO_HISTORY_FOUND) + "\n",
			wantStdErr: "",
			wantErr:    false,
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
		{
			name: "positive testing (one jrp in the database file)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShowRunE(nil, nil); err != nil {
						t.Errorf("historyShowOption.historyShowRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 1\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShowRunE(nil, nil); err != nil {
						t.Errorf("historyShowOption.historyShowRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest1\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n2\ttest2\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 2\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
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
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (Args is nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShowRunE(nil, nil); err != nil {
						t.Errorf("historyShowOption.historyShowRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 1\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (Args is empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShowRunE(nil, nil); err != nil {
						t.Errorf("historyShowOption.historyShowRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 1\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (DBFileDirPathProvider.GetJrpDBFileDirPath() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: mockDBFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShowRunE(nil, nil); err != nil {
						fmt.Printf("historyShowOption.historyShowRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: testutility.TEST_OUTPUT_ANY,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (Utility.CreateDirIfNotExitst() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               mockUtility,
					}
					if err := historyShowOption.historyShowRunE(nil, nil); err != nil {
						fmt.Printf("historyShowOption.historyShowRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: testutility.TEST_OUTPUT_ANY,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
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
				t.Errorf("historyShowOption.historyShowRunE() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("historyShowOption.historyShowRunE() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_historyShowOption_historyShow(t *testing.T) {
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
	fmtProxy := fmtproxy.New()
	jrpRepository := repository.New(
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
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME)
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
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing (no jrps in the database file)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_MESSAGE_NO_HISTORY_FOUND) + "\n",
			wantStdErr: "",
			wantErr:    false,
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
		{
			name: "positive testing (one jrp in the database file)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 1\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest1\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n2\ttest2\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 2\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
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
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (11 jrps in the database file, all)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   true,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest1\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n2\ttest2\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n3\ttest3\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n4\ttest4\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n5\ttest5\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n6\ttest6\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n7\ttest7\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n8\ttest8\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n9\ttest9\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n10\ttest10\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n11\ttest11\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\nTOTAL : 11\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
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
						{
							Phrase:    "test3",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test4",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test5",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test6",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test7",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test8",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test9",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test10",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test11",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (11 jrps in the database file, not all)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                0,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n2\ttest2\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n3\ttest3\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n4\ttest4\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n5\ttest5\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n6\ttest6\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n7\ttest7\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n8\ttest8\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n9\ttest9\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n10\ttest10\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n11\ttest11\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\nTOTAL : 10\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
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
						{
							Phrase:    "test3",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test4",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test5",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test6",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test7",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test8",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test9",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test10",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test11",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (11 jrps in the database file, not all, number is 9)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                9,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n3\ttest3\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n4\ttest4\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n5\ttest5\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n6\ttest6\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n7\ttest7\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n8\ttest8\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n9\ttest9\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n10\ttest10\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n11\ttest11\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\nTOTAL : 9\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
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
						{
							Phrase:    "test3",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test4",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test5",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test6",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test7",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test8",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test9",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test10",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test11",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (11 jrps in the database file, not all, number is 10)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                10,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n2\ttest2\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n3\ttest3\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n4\ttest4\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n5\ttest5\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n6\ttest6\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n7\ttest7\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n8\ttest8\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n9\ttest9\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n10\ttest10\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n11\ttest11\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\nTOTAL : 10\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
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
						{
							Phrase:    "test3",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test4",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test5",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test6",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test7",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test8",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test9",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test10",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test11",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (11 jrps in the database file, not all, number is 11)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                11,
						All:                   false,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					if err := historyShowOption.historyShow(jrpDBFilePath); err != nil {
						t.Errorf("historyShowOption.historyShow() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest1\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n2\ttest2\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n3\ttest3\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n4\ttest4\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n5\ttest5\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n6\ttest6\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n7\ttest7\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n8\ttest8\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n9\ttest9\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n10\ttest10\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n11\ttest11\t\t\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\nTOTAL : 11\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
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
						{
							Phrase:    "test3",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test4",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test5",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test6",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test7",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test8",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test9",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test10",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test11",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					},
				); err != nil {
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
				t.Errorf("historyShowOption.historyShow() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("historyShowOption.historyShow() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_historyShowOption_writeHistoryShowResult(t *testing.T) {
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
	fmtProxy := fmtproxy.New()
	jrpRepository := repository.New(
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
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult(nil)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_SHOW_MESSAGE_NO_HISTORY_FOUND) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are empty, not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult([]model.Jrp{})
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_SHOW_MESSAGE_NO_HISTORY_FOUND) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps is one, not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult([]model.Jrp{
						{
							ID:          1,
							Phrase:      "test",
							Prefix:      sqlProxy.StringToNullString("prefix"),
							Suffix:      sqlProxy.StringToNullString("suffix"),
							IsFavorited: 0,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					})
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest\tprefix\tsuffix\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 1\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps is two, not plain)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult([]model.Jrp{
						{
							ID:          1,
							Phrase:      "test1",
							Prefix:      sqlProxy.StringToNullString("prefix1"),
							Suffix:      sqlProxy.StringToNullString("suffix1"),
							IsFavorited: 0,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							ID:          2,
							Phrase:      "test2",
							Prefix:      sqlProxy.StringToNullString("prefix2"),
							Suffix:      sqlProxy.StringToNullString("suffix2"),
							IsFavorited: 1,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					})
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest1\tprefix1\tsuffix1\t\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n2\ttest2\tprefix2\tsuffix2\t○\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 2\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are nil, plain)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult(nil)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_SHOW_MESSAGE_NO_HISTORY_FOUND) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are empty, plain)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult([]model.Jrp{})
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_MESSAGE_NO_HISTORY_FOUND) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps is one, plain)",
			fields: fields{
				t: t,
				fnc: func() {
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult([]model.Jrp{
						{
							ID:          1,
							Phrase:      "test",
							Prefix:      sqlProxy.StringToNullString("prefix"),
							Suffix:      sqlProxy.StringToNullString("suffix"),
							IsFavorited: 0,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
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
					historyShowOption := &historyShowOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Number:                1,
						Plain:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						JrpWriter:             jrpWriter,
						Utility:               util,
					}
					historyShowOption.writeHistoryShowResult([]model.Jrp{
						{
							ID:          1,
							Phrase:      "test1",
							Prefix:      sqlProxy.StringToNullString("prefix1"),
							Suffix:      sqlProxy.StringToNullString("suffix1"),
							IsFavorited: 0,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							ID:          2,
							Phrase:      "test2",
							Prefix:      sqlProxy.StringToNullString("prefix2"),
							Suffix:      sqlProxy.StringToNullString("suffix2"),
							IsFavorited: 1,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
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
				t.Errorf("historyShowOption.writeHistoryShowResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("historyShowOption.writeHistoryShowResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
