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
	"github.com/yanosea/jrp/app/proxy/promptui"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/time"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"

	"github.com/yanosea/jrp/mock/app/database/jrp/repository"
	"github.com/yanosea/jrp/mock/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/mock/app/library/utility"
	"github.com/yanosea/jrp/mock/app/proxy/promptui"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNewFavoriteClearCommand(t *testing.T) {
	type args struct {
		g             *GlobalOption
		promptuiProxy promptuiproxy.Promptui
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		setup     func(*gomock.Controller, *args)
	}{
		{
			name: "positive testing",
			args: args{
				g: NewGlobalOption(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
				promptuiProxy: nil,
			},
			wantError: false,
			setup: func(mockCtrl *gomock.Controller, args *args) {
				mockPromptInstance := mockpromptuiproxy.NewMockPromptInstanceInterface(mockCtrl)
				mockPromptInstance.EXPECT().SetLabel(gomock.Any())
				mockPromptInstance.EXPECT().Run().Return("y", nil)
				mockPromptuiProxy := mockpromptuiproxy.NewMockPromptui(mockCtrl)
				mockPromptuiProxy.EXPECT().NewPrompt().Return(mockPromptInstance)
				args.promptuiProxy = mockPromptuiProxy
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
			got := NewFavoriteClearCommand(tt.args.g, tt.args.promptuiProxy)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewFavoriteClearCommand().Execute() : error =\n%v", err)
			}
		})
	}
}

func Test_favoriteClearOption_favoriteClearRunE(t *testing.T) {
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
	mockPromptInstance := mockpromptuiproxy.NewMockPromptInstanceInterface(mockCtrl)
	mockPromptInstance.EXPECT().SetLabel(gomock.Any()).AnyTimes()
	mockPromptInstance.EXPECT().Run().Return("y", nil).AnyTimes()
	mockPromptuiProxy := mockpromptuiproxy.NewMockPromptui(mockCtrl)
	mockPromptuiProxy.EXPECT().NewPrompt().Return(mockPromptInstance).AnyTimes()
	noMockPromptInstance := mockpromptuiproxy.NewMockPromptInstanceInterface(mockCtrl)
	noMockPromptInstance.EXPECT().SetLabel(gomock.Any())
	noMockPromptInstance.EXPECT().Run().Return("n", nil)
	noMockPromptuiProxy := mockpromptuiproxy.NewMockPromptui(mockCtrl)
	noMockPromptuiProxy.EXPECT().NewPrompt().Return(noMockPromptInstance)
	errMockPromptInstance := mockpromptuiproxy.NewMockPromptInstanceInterface(mockCtrl)
	errMockPromptInstance.EXPECT().SetLabel(gomock.Any())
	errMockPromptInstance.EXPECT().Run().Return("", errors.New("PromptInstance.Run() failed"))
	errMockPromptuiProxy := mockpromptuiproxy.NewMockPromptui(mockCtrl)
	errMockPromptuiProxy.EXPECT().NewPrompt().Return(errMockPromptInstance)
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
			name: "positive testing (answering not yes on the prompt)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         noMockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: constant.FAVORITE_CLEAR_MESSAGE_CLEAR_CANCELED + "\n",
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
			name: "positive testing (set no-confirm option)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             true,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         nil,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_NONE) + "\n",
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
			name: "positive testing (no jrps in the database file)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_NONE) + "\n",
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
			name: "positive testing (one jrp in the database file, not favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_NONE) + "\n",
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
			name: "positive testing (one jrp in the database file, favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY) + "\n",
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
			name: "positive testing (two jrps in the database file, both are not favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_NONE) + "\n",
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
			name: "positive testing (two jrps in the database file, both are favorited)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY) + "\n",
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
			name: "positive testing (two jrps in the database file, one is favorited and other is not)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						t.Errorf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
					}
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY) + "\n",
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
			name: "negative testing (DBFileDirPathProvider.GetJrpDBFileDirPath() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: mockDBFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						fmt.Printf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
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
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               mockUtility,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						fmt.Printf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
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
			name: "negative testing (PromptInstance.run() failed)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         errMockPromptuiProxy,
						Utility:               util,
					}
					if err := favoriteClearOption.favoriteClearRunE(nil, nil); err != nil {
						fmt.Printf("favoriteClearOption.favoriteClearRunE() : error =\n%v", err)
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
				t.Errorf("favoriteClearOption.favoriteClearRunE() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != testutility.TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("favoriteClearOption.favoriteClearRunE() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(favoritedJrps, tt.wantJrps) {
				t.Errorf("favoriteClearOption.favoriteClearRunE() : favoritedJrps =\n%v, want =\n%v", favoritedJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_favoriteClearOption_favoriteClear(t *testing.T) {
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPromptInstance := mockpromptuiproxy.NewMockPromptInstanceInterface(mockCtrl)
	mockPromptInstance.EXPECT().SetLabel(gomock.Any()).AnyTimes()
	mockPromptInstance.EXPECT().Run().Return("y", nil).AnyTimes()
	mockPromptuiProxy := mockpromptuiproxy.NewMockPromptui(mockCtrl)
	mockPromptuiProxy.EXPECT().NewPrompt().Return(mockPromptInstance).AnyTimes()

	type fields struct {
		Out                   ioproxy.WriterInstanceInterface
		ErrOut                ioproxy.WriterInstanceInterface
		Args                  []string
		NoConfirm             bool
		DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
		JrpRepository         repository.JrpRepositoryInterface
		PromptuiProxy         promptuiproxy.Promptui
		Utility               utility.UtilityInterface
	}
	type args struct {
		jrpDBFilePath string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantJrps []*model.Jrp
		wantErr  bool
		setup    func(*gomock.Controller, *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no jrps in the database file)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				NoConfirm:             false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
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
			name: "positive testing (one jrp in the database file, not favorited)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				NoConfirm:             false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (one jrp in the database file, favorited)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				NoConfirm:             false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, both are not favorited)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				NoConfirm:             false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, both are favorited)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				NoConfirm:             false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "positive testing (two jrps in the database file, one is favorited and other is not)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				NoConfirm:             false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         jrpRepository,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			wantJrps: nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, _ *fields) {
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
			name: "negative testing (JrpRepository.RemoveFavoriteAll() failed)",
			fields: fields{
				Out:                   osproxy.Stdout,
				ErrOut:                osproxy.Stderr,
				Args:                  osproxy.Args[1:],
				NoConfirm:             false,
				DBFileDirPathProvider: dbFileDirPathProvider,
				JrpRepository:         nil,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               util,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			wantJrps: nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
				}
				mockJrpRepository := mockrepository.NewMockJrpRepositoryInterface(mockCtrl)
				mockJrpRepository.EXPECT().RemoveFavoriteAll(gomock.Any()).Return(repository.RemovedFailed, errors.New("JrpRepository.RemoveFavoriteAll() : error"))
				tt.JrpRepository = mockJrpRepository
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
			o := &favoriteClearOption{
				Out:                   tt.fields.Out,
				ErrOut:                tt.fields.ErrOut,
				Args:                  tt.fields.Args,
				NoConfirm:             false,
				DBFileDirPathProvider: tt.fields.DBFileDirPathProvider,
				JrpRepository:         tt.fields.JrpRepository,
				PromptuiProxy:         mockPromptuiProxy,
				Utility:               tt.fields.Utility,
			}
			if err := o.favoriteClear(tt.args.jrpDBFilePath); (err != nil) != tt.wantErr {
				t.Errorf("favoriteClearOption.favoriteClear() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(favoritedJrps, tt.wantJrps) {
				t.Errorf("favoriteClearOption.favoriteClear() : favoritedJrps =\n%v, want =\n%v", favoritedJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func Test_favoriteClearOption_writeFavoriteClearResult(t *testing.T) {
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPromptInstance := mockpromptuiproxy.NewMockPromptInstanceInterface(mockCtrl)
	mockPromptInstance.EXPECT().SetLabel(gomock.Any()).AnyTimes()
	mockPromptInstance.EXPECT().Run().Return("y", nil).AnyTimes()
	mockPromptuiProxy := mockpromptuiproxy.NewMockPromptui(mockCtrl)
	mockPromptuiProxy.EXPECT().NewPrompt().Return(mockPromptInstance).AnyTimes()

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
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					favoriteClearOption.writeFavoriteClearResult(repository.RemovedSuccessfully)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.GreenString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY) + "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (result is RemovedFailed)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					favoriteClearOption.writeFavoriteClearResult(repository.RemovedFailed)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: colorProxy.RedString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_FAIRULE) + "\n",
			wantErr:    false,
		},
		{
			name: "positive testing (result is RemovedNone)",
			fields: fields{
				t: t,
				fnc: func() {
					favoriteClearOption := &favoriteClearOption{
						Out:                   capturer.OutBuffer,
						ErrOut:                capturer.ErrBuffer,
						Args:                  osproxy.Args[1:],
						NoConfirm:             false,
						DBFileDirPathProvider: dbFileDirPathProvider,
						JrpRepository:         jrpRepository,
						PromptuiProxy:         mockPromptuiProxy,
						Utility:               util,
					}
					favoriteClearOption.writeFavoriteClearResult(repository.RemovedNone)
				},
				capturer: capturer,
			},
			wantStdOut: colorProxy.YellowString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_NONE) + "\n",
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
				t.Errorf("favoriteClearOption.writeFavoriteClearResult() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("favoriteClearOption.writeFavoriteClearResult() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
