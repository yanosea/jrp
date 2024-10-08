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

func TestNewFavoriteAddCommand(t *testing.T) {
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
			got := NewFavoriteAddCommand(tt.args.g)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewFavoriteAddCommand().Execute() : error =\n%v", err)
			}
		})
	}
}

func Test_favoriteAddOption_favoriteAddRunE(t *testing.T) {
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
		wantJrps   []*model.Jrp
		wantErr    bool
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing (no jrps in the database file, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
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
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
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
			name: "positive testing (one jrp in the database file, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, not favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, not favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, not favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, not favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, not favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, not favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				}, {
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
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (DBFileDirPathProvider.GetJrpDBFileDirPath() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "remove", "1", "2"},
						DBFileDirPathProvider: mockDBFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						fmt.Printf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
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
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (Utility.CreateDirIfNotExitst() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "remove", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               mockUtility,
					}
					if err := favoriteAddOption.favoriteAddRunE(nil, nil); err != nil {
						fmt.Printf("favoriteAddOption.favoriteAddRunE() : error =\n%v", err)
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
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
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
				t.Errorf("favoriteAddOption.favoriteAddRunE() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("favoriteAddOption.favoriteAddRunE() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(favoritedJrps, tt.wantJrps) {
				t.Errorf("favoriteAddOption.favoriteAddRunE() : favoritedJrps =\n%v, want =\n%v", favoritedJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_favoriteAddOption_favoriteAdd(t *testing.T) {
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
		wantJrps   []*model.Jrp
		wantErr    bool
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing (no jrps in the database file, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
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
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
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
			name: "positive testing (one jrp in the database file, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, not favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{2}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, not favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{1}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, not favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{1, 2}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, not favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{3}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps:   nil,
			wantErr:    false,
			setup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, not favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{1, 2}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, not favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{2, 3}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NOT_ALL) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{2}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{1}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (one jrp in the database file, favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{1, 2}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, args are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  nil,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, nil); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, args are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, id in args does not match any jrps)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{3}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, id in args matches a jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "1", "2"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{1, 2}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
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
					[]*model.Jrp{
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
			name: "positive testing (two jrps in the database file, favorited, one of two ids in args matches jrps but other does not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  []string{"favorite", "add", "2", "3"},
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					if err := favoriteAddOption.favoriteAdd(jrpDBFilePath, []int{2, 3}); err != nil {
						t.Errorf("favoriteAddOption.favoriteAdd() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantJrps: []*model.Jrp{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				}, {
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
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]*model.Jrp{
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
					t.Errorf("OsProxy.AddAll() : error =\n%v", err)
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
				t.Errorf("favoriteAddOption.favoriteAdd() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("favoriteAddOption.favoriteAdd() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(favoritedJrps, tt.wantJrps) {
				t.Errorf("favoriteAddOption.favoriteAdd() : favoritedJrps =\n%v, want =\n%v", favoritedJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_favoriteAddOption_writeFavoriteAddResult(t *testing.T) {
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
			name: "positive testing (result is AddedSuccessfully)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					favoriteAddOption.writeFavoriteAddResult(repository.AddedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is AddedFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					favoriteAddOption.writeFavoriteAddResult(repository.AddedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.FAVORITE_ADD_MESSAGE_ADDED_FAILURE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (result is AddedNone)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					favoriteAddOption.writeFavoriteAddResult(repository.AddedNone)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is AddedNotAll)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteAddOption := &favoriteAddOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						Utility:               util,
					}
					favoriteAddOption.writeFavoriteAddResult(repository.AddedNotAll)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NOT_ALL) + "\n",
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
				t.Errorf("favoriteAddOption.writeFavoriteAddResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("favoriteAddOption.writeFavoriteAddResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
