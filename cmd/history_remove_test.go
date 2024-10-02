package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/database/jrp/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/time"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"

	"github.com/yanosea/jrp/mock/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/mock/app/library/utility"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNewHistoryRemoveCommand(t *testing.T) {
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
			got := NewHistoryRemoveCommand(tt.args.g)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewHistoryRemoveCommand().Execute() : error =\n%v", err)
			}
		})
	}
}

func Test_historyRemoveOption_historyRemoveRunE(t *testing.T) {
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
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	filepathProxy := filepathproxy.New()
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME)
	timeProxy := timeproxy.New()
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)
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
		wantJrps   []model.Jrp
		wantErr    bool
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing (no jrps in the database file, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
			name: "positive testing (no jrps in the database file, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
			name: "positive testing (one jrp in the database file, not favorited, not force, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
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
			wantErr: false,
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
			name: "positive testing (one jrp in the database file, not favorited, not force, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
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
			wantErr: false,
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
			name: "positive testing (one jrp in the database file, not favorited, not force, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "2"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
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
			wantErr: false,
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
			name: "positive testing (one jrp in the database file, not favorited, not force, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
			name: "positive testing (one jrp in the database file, not favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
			name: "positive testing (two jrps in the database file, not favorited, not force, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
			name: "positive testing (two jrps in the database file, not favorited, not force, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
			name: "positive testing (two jrps in the database file, not favorited, not force, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "3", "4"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
			name: "positive testing (two jrps in the database file, not favorited, not force, ids in args matches a jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
			name: "positive testing (two jrps in the database file, not favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "2", "3"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
			name: "positive testing (one jrp in the database file, favorited, not force, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "2"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "3", "4"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, ids in args matches jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "2", "3"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "2"},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1"},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "3", "4"},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, ids in args matches jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "2", "3"},
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						t.Errorf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
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
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 false,
						DBFileDirPathProvider: mockDBFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						fmt.Printf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: testutility.TEST_OUTPUT_ANY,
			wantJrps:   nil,
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
			name: "negative testing (Utility.CreateDirIfNotExitst() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"history", "remove", "1", "2"},
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               mockUtility,
					}
					if err := historyRemoveOption.historyRemoveRunE(nil, nil); err != nil {
						fmt.Printf("historyRemoveOption.historyRemoveRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: testutility.TEST_OUTPUT_ANY,
			wantJrps:   nil,
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
				t.Errorf("historyRemoveOption.historyRemoveRunE() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("historyRemoveOption.historyRemoveRunE() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			leftJrps, err := jrpRepository.GetAllHistory(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(leftJrps, tt.wantJrps) {
				t.Errorf("historyRemoveOption.historyRemoveRunE() : leftJrps =\n%v, want =\n%v", leftJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_historyRemoveOption_historyRemove(t *testing.T) {
	osProxy := osproxy.New()
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
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	filepathProxy := filepathproxy.New()
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME)
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
		Force                 bool
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		JrpRepository         repository.JrpRepositoryInterface
		Utility               utility.UtilityInterface
	}
	type args struct {
		jrpDBFilePath string
		IDs           []int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantJrps []model.Jrp
		wantErr  bool
		setup    func(*gomock.Controller, *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no jrps in the database file, ids are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  nil,
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           nil,
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
			name: "positive testing (no jrps in the database file, ids are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{},
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
			name: "positive testing (one jrp in the database file, not favorited, not force, ids are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  nil,
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           nil,
			},
			wantJrps: []model.Jrp{
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
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (one jrp in the database file, not favorited, not force, ids are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{},
			},
			wantJrps: []model.Jrp{
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
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (one jrp in the database file, not favorited, not force, id in args does not match any jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "2"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{2},
			},
			wantJrps: []model.Jrp{
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
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (one jrp in the database file, not favorited, not force, id in args matches a jrp)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (one jrp in the database file, not favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1", "2"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1, 2},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, not favorited, not force, ids are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  nil,
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           nil,
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, not favorited, not force, ids are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, not favorited, not force, id in args does not match any jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "3", "4"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{3, 4},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, not favorited, not force, ids in args matches a jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1", "2"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1, 2},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, not favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "2", "3"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{2, 3},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 0,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (one jrp in the database file, favorited, not force, ids are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  nil,
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           nil,
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, ids are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, id in args does not match any jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "2"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{2},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, id in args matches a jrp)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1", "2"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1, 2},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, ids are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  nil,
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           nil,
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, ids are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, id in args does not match any jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "3", "4"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{3, 4},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, ids in args matches jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1", "2"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1, 2},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, not force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "2", "3"},
				Force:                 false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{2, 3},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, ids are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  nil,
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           nil,
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, ids are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, id in args does not match any jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "2"},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{2},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, id in args matches a jrp)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1"},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (one jrp in the database file, favorited, force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1", "2"},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1, 2},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, ids are nil)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  nil,
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           nil,
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, ids are empty)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, id in args does not match any jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "3", "4"},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{3, 4},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, ids in args matches jrps)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "1", "2"},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{1, 2},
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two jrps in the database file, favorited, force, one of two ids in args matches jrps but other does not)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  []string{"history", "remove", "2", "3"},
				Force:                 true,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				IDs:           []int{2, 3},
			},
			wantJrps: []model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
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
			o := &historyRemoveOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Args:                  tt.fields.Args,
				Force:                 tt.fields.Force,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				JrpRepository:         tt.fields.JrpRepository,
				Utility:               tt.fields.Utility,
			}
			if err := o.historyRemove(tt.args.jrpDBFilePath, tt.args.IDs); (err != nil) != tt.wantErr {
				t.Errorf("historyRemoveOption.historyRemove() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			leftJrps, err := jrpRepository.GetAllHistory(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(leftJrps, tt.wantJrps) {
				t.Errorf("historyRemoveOption.historyRemove() : leftJrps =\n%v, want =\n%v", leftJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_historyRemoveOption_writeHistoryRemoveResult(t *testing.T) {
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
	util := utility.New(
		fmtProxy,
		osProxy,
		strconvproxy.New(),
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
			name: "positive testing (result is RemovedSuccessfully)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					historyRemoveOption.writeHistoryRemoveResult(repository.RemovedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is RemovedFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					historyRemoveOption.writeHistoryRemoveResult(repository.RemovedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (result is RemovedNone)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Force:                 false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					historyRemoveOption.writeHistoryRemoveResult(repository.RemovedNone)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is RemovedNotAll)",
			fields: fields{
				t: t,
				fnc: func() {
					historyRemoveOption := &historyRemoveOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						Force:                 true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					historyRemoveOption.writeHistoryRemoveResult(repository.RemovedNotAll)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL) + "\n",
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
				t.Errorf("historyRemoveOption.writeHistoryRemoveResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("historyRemoveOption.writeHistoryRemoveResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
