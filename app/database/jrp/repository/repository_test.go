package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/time"
	"github.com/yanosea/jrp/app/proxy/user"

	"github.com/yanosea/jrp/mock/app/proxy/sort"
	"github.com/yanosea/jrp/mock/app/proxy/sql"
	"github.com/yanosea/jrp/test/testutility"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	fmtProxy := fmtproxy.New()
	sortProxy := sortproxy.New()
	sqlProxy := sqlproxy.New()
	stringProxy := stringsproxy.New()

	type args struct {
		fmtProxy     fmtproxy.Fmt
		sortProxy    sortproxy.Sort
		sqlProxy     sqlproxy.Sql
		stringsProxy stringsproxy.Strings
	}
	tests := []struct {
		name string
		args args
		want *JrpRepository
	}{
		{
			name: "positive testing",
			args: args{
				fmtProxy:     fmtProxy,
				sortProxy:    sortProxy,
				sqlProxy:     sqlProxy,
				stringsProxy: stringProxy,
			},
			want: &JrpRepository{
				FmtProxy:     fmtProxy,
				SortProxy:    sortProxy,
				SqlProxy:     sqlProxy,
				StringsProxy: stringProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.fmtProxy, tt.args.sortProxy, tt.args.sqlProxy, tt.args.stringsProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want=\n%v", got, tt.want)
			}
		})
	}
}

func TestJrpRepository_SaveHistory(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		jrps          []model.Jrp
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus SaveStatus
		wantJrps   []model.Jrp
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller, tt *fields)
		cleanup    func()
	}{
		{
			name: "positive testing (jrps are nil)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps:          nil,
			},
			wantStatus: SavedNone,
			wantJrps:   nil,
			wantErr:    false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (jrps are empty)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps:          []model.Jrp{},
			},
			wantStatus: SavedNone,
			wantJrps:   nil,
			wantErr:    false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (jrps are one, there is no jrps in the database)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
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
			wantStatus: SavedSuccessfully,
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (jrps are two, there is no jrps in the database)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
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
			wantStatus: SavedSuccessfully,
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (jrps are one, there is jrps in the database alresdy)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps: []model.Jrp{
					{
						Phrase:    "test2",
						Prefix:    sqlProxy.StringToNullString(""),
						Suffix:    sqlProxy.StringToNullString(""),
						CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
			},
			wantStatus: SavedSuccessfully,
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							ID:        1,
							Phrase:    "test1",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (jrps are two, there is jrps in the databasea already)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				jrps: []model.Jrp{
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
				},
			},
			wantStatus: SavedSuccessfully,
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
				{
					ID:        3,
					Phrase:    "test3",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							ID:        1,
							Phrase:    "test1",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedFailed,
			wantJrps:   nil,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedFailed,
			wantJrps:   nil,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Begin() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedFailed,
			wantJrps:   nil,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(nil, errors.New("DBInstance.Begin() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Prepare() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedFailed,
			wantJrps:   nil,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(nil, errors.New("DBInstance.Prepare() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (stmt.Exec() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedFailed,
			wantJrps:   nil,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("StmtInstance.Exec() failed"))
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (res.RowsAffected() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedFailed,
			wantJrps:   nil,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(0), errors.New("ResultInstance.RowsAffected() failed"))
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (tx.Commit() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedFailed,
			wantJrps:   nil,
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(0), nil)
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Commit().Return(errors.New("TxInstance.Commit() failed"))
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (res.RowsAffected() != len(jrps))",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
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
			wantStatus: SavedNotAll,
			wantJrps:   nil,
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(10), nil)
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Commit().Return(nil)
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.SaveHistory(tt.args.jrpDBFilePath, tt.args.jrps)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.SaveHistory() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if got != tt.wantStatus {
				t.Errorf("JrpRepository.SaveHistory() : got =\n%v, want =\n%v", got, tt.wantStatus)
			}
			savedJrps, err := jrpRepository.GetAllHistory(tt.args.jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if !jrpChecker.IsSameJrps(savedJrps, tt.wantJrps) {
				t.Errorf("JrpRepository.SaveHistory() : savedJrps =\n%v, want =\n%v", savedJrps, tt.wantJrps)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_GetAllHistory(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy      fmtproxy.Fmt
		SortProxy     sortproxy.Sort
		SqlProxy      sqlproxy.Sql
		StringsProxy  stringsproxy.Strings
		JrpDBFilePath string
	}
	type args struct {
		jrpDBFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (there are no jrps in the database)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is two jrps in the database)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any()).Return(nil, errors.New("DBInstance.Query() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any()).Return(mockRowsInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.GetAllHistory(tt.args.jrpDBFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.GetAllHistory() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_GetHistoryWithNumber(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		number        int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (there are no jrps in the database, number is -1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        -1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are no jrps in the database, number is 0)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are no jrps in the database, number is 1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are no jrps in the database, number is 2)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is -1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        -1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 0)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 2)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is -1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        -1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 0)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 2)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Prepare() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(nil, errors.New("DBInstance.Prepare() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (stmt.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Query(gomock.Any()).Return(nil, errors.New("StmtInstance.Query() failed"))
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Query(gomock.Any()).Return(mockRowsInstance, nil)
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SortProxy.Slice() returns false)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want: []model.Jrp{
				{
					ID:        2,
					Phrase:    "test2",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
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
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				mockSortProxy := mocksortproxy.NewMockSort(mockCtrl)
				mockSortProxy.EXPECT().Slice(gomock.Any(), gomock.Any())
				tt.SortProxy = mockSortProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.GetHistoryWithNumber(tt.args.jrpDBFilePath, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.GetHistoryWithNumber() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.GetHistoryWithNumber() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_SearchAllHistory(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		keywords      []string
		and           bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (keywords are nil, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      nil,
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (keywords are nil, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      nil,
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (keywords are empty, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (keywords are empty, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, do not match any jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"テスト"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, do not match any jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"テスト"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match one jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match one jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match one jrp, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match one jrp, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛として時雨",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match two jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:        2,
					Phrase:    "時雨",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match two jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛として時雨1",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:        2,
					Phrase:    "凛として時雨2",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				keywords: []string{"test"},
				and:      false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("DBInstance.Query() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(mockRowsInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.SearchAllHistory(tt.args.jrpDBFilePath, tt.args.keywords, tt.args.and)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.SearchAllHistory() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.SearchAllHistory() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_SearchHistoryWithNumber(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		number        int
		keywords      []string
		and           bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(*gomock.Controller, *fields)
		cleanup func()
	}{
		{
			name: "positive testing (number is 0, keywords are nil, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      nil,
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, keywords are nil, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      nil,
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0 keywords are nil, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      nil,
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, keywords are empty, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, keywords are empty, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, single keyword, do not match any jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{"テスト"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, single keyword, do not match any jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{"テスト"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match one jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match one jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match one jrp, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match one jrp, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛として時雨",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:        2,
					Phrase:    "時雨",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        2,
					Phrase:    "時雨",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        1,
					Phrase:    "凛として時雨1",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:        2,
					Phrase:    "凛として時雨2",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:        2,
					Phrase:    "凛として時雨2",
					Prefix:    sqlProxy.StringToNullString(""),
					Suffix:    sqlProxy.StringToNullString(""),
					CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("DBInstance.Query() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(mockRowsInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.SearchHistoryWithNumber(tt.args.jrpDBFilePath, tt.args.number, tt.args.keywords, tt.args.and)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.SearchHistoryWithNumber() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.SearchHistoryWithNumber() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_RemoveHistoryByIDs(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		ids           []int
		force         bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RemoveStatus
		wantErr bool
		setup   func(*gomock.Controller, *fields)
		cleanup func()
	}{
		{
			name: "positive testing (ids are nil, not force)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           nil,
				force:         false,
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are empty, not force)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{},
				force:         false,
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, not force, database does not have the id)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{2},
				force:         false,
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, not force, database has the id and it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         false,
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, not force, database has the id and it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         false,
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, force, database has the id and it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         true,
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, force, database has the id and it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         true,
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, not force, database does not have the id)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{2, 3},
				force:         false,
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, not force, database has the id and both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
				force:         false,
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, not force, database has the id and both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
				force:         false,
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, not force, database has the id and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
				force:         false,
			},
			want:    RemovedNotAll,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, force, database has the id and both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
				force:         true,
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, force, database has the id and both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
				force:         true,
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, force, database has the id and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
				force:         true,
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         false,
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         false,
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Prepare() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         false,
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(nil, errors.New("DBInstance.Prepare() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (stmt.Exec() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         false,
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("StmtInstance.Exec() failed"))
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (res.RowsAffected() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
				force:         false,
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(0), errors.New("ResultInstance.RowsAffected() failed"))
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.RemoveHistoryByIDs(tt.args.jrpDBFilePath, tt.args.ids, tt.args.force)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.RemoveHistoryByIDs() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("JrpRepository.RemoveHistoryByIDs() : got =\n%v, want =\n%v", got, tt.want)
			}
			if err == nil && tt.args.ids != nil && len(tt.args.ids) > 0 {
				for _, id := range tt.args.ids {
					isExist, err := jrpChecker.IsExist(jrpDBFilePath, id)
					if err != nil {
						t.Errorf("JrpChecker.IsExist() : error =\n%v", err)
					}
					if isExist {
						isFavorited, err := jrpChecker.IsFavorited(jrpDBFilePath, id)
						if err != nil {
							t.Errorf("JrpChecker.IsFavorited() : error =\n%v", err)
						}
						if !isFavorited || tt.args.force {
							t.Errorf("JrpRepository.RemoveHistoryByIDs() : did not removed \n[%v]", id)
						}
					}
				}
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_RemoveHistoryAll(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		force         bool
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		want             RemoveStatus
		wantLeftJrpCount int
		wantSeq          int
		wantErr          bool
		setup            func(*gomock.Controller, *fields)
		cleanup          func()
	}{
		{
			name: "positive testing (not force, there is no jrps in the database)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedNone,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not force, there is one jrp in the database and it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not force, there is one jrp in the database and it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedNone,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not force, there is two jrps in the database and both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not force, there is two jrps in the database and both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedNone,
			wantLeftJrpCount: 2,
			wantSeq:          2,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (not force, there is two jrps in the database and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 1,
			wantSeq:          2,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (force, there is no jrps in the database)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         true,
			},
			want:             RemovedNone,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (force, there is one jrp in the database and it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         true,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (force, there is one jrp in the database and it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         true,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (force, there is two jrps in the database and both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         true,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (force, there is two jrps in the database and both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         true,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (force, there is two jrps in the database and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         true,
			},
			want:             RemovedSuccessfully,
			wantLeftJrpCount: 0,
			wantSeq:          0,
			wantErr:          false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Begin() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(nil, errors.New("DBInstance.Begin() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (tx.Exec(q) failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("TxInstance.Exec() failed"))
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (tx.QueryRow().Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowInstance := mocksqlproxy.NewMockRowInstanceInterface(mockCtrl)
				mockRowInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowInstance.Scan() failed"))
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockTxInstance.EXPECT().QueryRow(gomock.Any()).Return(mockRowInstance)
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (tx.Exec(query.RemoveJrpSeq) failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowInstance := mocksqlproxy.NewMockRowInstanceInterface(mockCtrl)
				mockRowInstance.EXPECT().Scan(gomock.Any()).Return(nil)
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockTxInstance.EXPECT().QueryRow(gomock.Any()).Return(mockRowInstance)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, errors.New("TxInstance.Exec() failed"))
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (res.RowsAffected() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowInstance := mocksqlproxy.NewMockRowInstanceInterface(mockCtrl)
				mockRowInstance.EXPECT().Scan(gomock.Any()).Return(nil)
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(0), errors.New("ResultInstance.RowsAffected() failed"))
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockTxInstance.EXPECT().QueryRow(gomock.Any()).Return(mockRowInstance)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (tx.Commit() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				force:         false,
			},
			want:             RemovedFailed,
			wantLeftJrpCount: 1,
			wantSeq:          1,
			wantErr:          true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowInstance := mocksqlproxy.NewMockRowInstanceInterface(mockCtrl)
				mockRowInstance.EXPECT().Scan(gomock.Any()).Return(nil)
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(1), nil)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockTxInstance.EXPECT().QueryRow(gomock.Any()).Return(mockRowInstance)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockTxInstance.EXPECT().Commit().Return(errors.New("TxInstance.Commit() failed"))
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.RemoveHistoryAll(tt.args.jrpDBFilePath, tt.args.force)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.RemoveHistoryAll() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("JrpRepository.RemoveHistoryAll() : got =\n%v, want =\n%v", got, tt.want)
			}
			leftJrps, err := jrpRepository.GetAllHistory(tt.args.jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if len(leftJrps) != tt.wantLeftJrpCount {
				t.Errorf("JrpRepository.RemoveHistoryAll() : len(leftJrps) =\n%v, want =\n%v", leftJrps, tt.wantLeftJrpCount)
			}
			seq, err := jrpChecker.GetJrpSeq(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpChecker.GetJrpSeq() : error =\n%v", err)
			}
			if seq != tt.wantSeq {
				t.Errorf("JrpRepository.RemoveHistoryAll() : seq =\n%v, want =\n%v", seq, tt.wantSeq)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_GetAllFavorite(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (there are no jrps in the database)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database and it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database and it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrp in the database and both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrp in the database and both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrp in the database and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("DBInstance.Query() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(mockRowsInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.GetAllFavorite(tt.args.jrpDBFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.GetAllFavorite() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_GetFavoriteWithNumber(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		number        int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (there are no jrps in the database, number is -1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        -1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are no jrps in the database, number is 0)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are no jrps in the database, number is 1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are no jrps in the database, number is 2)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is -1)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        -1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 0)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 1, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 1, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 2, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database, number is 2, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 0)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 1, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 1, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 2, both not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 2, both favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there are two jrps in the database, number is 2, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Prepare() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(nil, errors.New("DBInstance.Prepare() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (stmt.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Query(gomock.Any()).Return(nil, errors.New("StmtInstance.Query() failed"))
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Query(gomock.Any()).Return(mockRowsInstance, nil)
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SortProxy.Slice() returns false)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
			},
			want: []model.Jrp{
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
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
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				mockSortProxy := mocksortproxy.NewMockSort(mockCtrl)
				mockSortProxy.EXPECT().Slice(gomock.Any(), gomock.Any())
				tt.SortProxy = mockSortProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.GetFavoriteWithNumber(tt.args.jrpDBFilePath, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.GetFavoriteWithNumber() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.GetFavoriteWithNumber() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_SearchAllFavorite(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		keywords      []string
		and           bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (keywords are nil, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      nil,
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (keywords are nil, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      nil,
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (keywords are empty, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (keywords are empty, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, do not match any jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"テスト"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, do not match any jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"テスト"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match one jrps, OR condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match one jrps, OR condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match one jrps, AND condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match one jrps, AND condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, OR condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, OR condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, OR condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, AND condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, AND condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (single keyword, match two jrps, AND condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match one jrp, OR condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match one jrp, OR condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match one jrp, OR condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match one jrp, AND condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match one jrp, AND condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛として時雨",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match two jrps, OR condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match two jrps, OR condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "時雨",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match two jrps, AND condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match two jrps, AND condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛として時雨1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "凛として時雨2",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (two keywords, match two jrps, AND condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛として時雨1",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("DBInstance.Query() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(mockRowsInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.SearchAllFavorite(tt.args.jrpDBFilePath, tt.args.keywords, tt.args.and)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.SearchAllFavorite() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.SearchAllFavorite() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_SearchFavoriteWithNumber(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		number        int
		keywords      []string
		and           bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Jrp
		wantErr bool
		setup   func(*gomock.Controller, *fields)
		cleanup func()
	}{
		{
			name: "positive testing (number is 0, keywords are nil, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      nil,
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, keywords are nil, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      nil,
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0 keywords are nil, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      nil,
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, keywords are empty, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, keywords are empty, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, single keyword, do not match any jrps, OR condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{"テスト"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 0, single keyword, do not match any jrps, AND condition)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        0,
				keywords:      []string{"テスト"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match one jrps, OR condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match one jrps, OR condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match one jrps, AND condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match one jrps, AND condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, OR condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, OR condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, OR condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, OR condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, OR condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, OR condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, AND condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, AND condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, single keyword, match two jrps, AND condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, AND condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           true,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, AND condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, single keyword, match two jrps, AND condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           true,
			},
			want: []model.Jrp{
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match one jrp, OR condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match one jrp, OR condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match one jrp, AND condition, it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match one jrp, AND condition, it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛として時雨",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, OR condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, OR condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "時雨",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, OR condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, OR condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, OR condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          2,
					Phrase:      "時雨",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, OR condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "時雨",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, AND condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, AND condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛として時雨1",
					Prefix:      sqlProxy.StringToNullString(""),
					Suffix:      sqlProxy.StringToNullString(""),
					IsFavorited: 1,
					CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				},
				{
					ID:          2,
					Phrase:      "凛として時雨2",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 2, two keywords, match two jrps, AND condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        2,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛として時雨1",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, AND condition, both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, AND condition, both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          2,
					Phrase:      "凛として時雨2",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1, 2}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (number is 1, two keywords, match two jrps, AND condition, the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"凛", "時雨"},
				and:           false,
			},
			want: []model.Jrp{
				{
					ID:          1,
					Phrase:      "凛として時雨1",
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(
					jrpDBFilePath,
					[]model.Jrp{
						{
							Phrase:    "凛として時雨1",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "凛として時雨2",
							Prefix:    sqlProxy.StringToNullString(""),
							Suffix:    sqlProxy.StringToNullString(""),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("DBInstance.Query() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				number:        1,
				keywords:      []string{"test"},
				and:           false,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockRowsInstance := mocksqlproxy.NewMockRowsInstanceInterface(mockCtrl)
				mockRowsInstance.EXPECT().Next().Return(true)
				mockRowsInstance.EXPECT().Scan(gomock.Any()).Return(errors.New("RowsInstance.Scan() failed"))
				mockRowsInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Query(gomock.Any(), gomock.Any()).Return(mockRowsInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.SearchFavoriteWithNumber(tt.args.jrpDBFilePath, tt.args.number, tt.args.keywords, tt.args.and)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.SearchFavoriteWithNumber() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if !jrpChecker.IsSameJrps(got, tt.want) {
				t.Errorf("JrpRepository.SearchFavoriteWithNumber() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_AddFavoriteByIDs(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		ids           []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    AddStatus
		wantErr bool
		setup   func(*gomock.Controller, *fields)
		cleanup func()
	}{
		{
			name: "positive testing (ids are nil)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           nil,
			},
			want:    AddedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are empty)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{},
			},
			want:    AddedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, database does not have the id)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{2},
			},
			want:    AddedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, database has the id and it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    AddedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, database has the id and it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    AddedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, database does not have the id)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{2, 3},
			},
			want:    AddedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, database has the id and both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
			},
			want:    AddedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, database has the id and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
			},
			want:    AddedNotAll,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    AddedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    AddedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Prepare() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    AddedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(nil, errors.New("DBInstance.Prepare() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (stmt.Exec() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    AddedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("StmtInstance.Exec() failed"))
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (res.RowsAffected() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    AddedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(0), errors.New("ResultInstance.RowsAffected() failed"))
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.AddFavoriteByIDs(tt.args.jrpDBFilePath, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("JrpRepository.AddFavoriteByIDs() : got =\n%v, want =\n%v", got, tt.want)
			}
			if err == nil && tt.args.ids != nil && len(tt.args.ids) > 0 {
				for _, id := range tt.args.ids {
					isExist, err := jrpChecker.IsExist(jrpDBFilePath, id)
					if err != nil {
						t.Errorf("JrpChecker.IsExist() : error =\n%v", err)
					}
					if isExist {
						isFavorited, err := jrpChecker.IsFavorited(jrpDBFilePath, id)
						if err != nil {
							t.Errorf("JrpChecker.IsFavorited() : error =\n%v", err)
						}
						if !isFavorited {
							t.Errorf("JrpRepository.AddFavoriteByIDs() : did not favorited \n[%v]", id)
						}
					}
				}
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_RemoveFavoriteByIDs(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	jrpChecker := testutility.NewJrpChecker(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		strconvproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		ids           []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RemoveStatus
		wantErr bool
		setup   func(*gomock.Controller, *fields)
		cleanup func()
	}{
		{
			name: "positive testing (ids are nil)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           nil,
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are empty)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{},
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, database does not have the id)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{2},
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, database has the id and it is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids is one, database has the id and it is favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, database does not have the id)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{2, 3},
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, database has the id and both are not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
			},
			want:    RemovedNone,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, not force, database has the id and both are favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
			},
			want:    RemovedSuccessfully,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (ids are two, database has the id and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1, 2},
			},
			want:    RemovedNotAll,
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Prepare() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(nil, errors.New("DBInstance.Prepare() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (stmt.Exec() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("StmtInstance.Exec() failed"))
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (res.RowsAffected() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				ids:           []int{1},
			},
			want:    RemovedFailed,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(0), errors.New("ResultInstance.RowsAffected() failed"))
				mockStmtInstance := mocksqlproxy.NewMockStmtInstanceInterface(mockCtrl)
				mockStmtInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockStmtInstance.EXPECT().Close().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Prepare(gomock.Any()).Return(mockStmtInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.RemoveFavoriteByIDs(tt.args.jrpDBFilePath, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.RemoveFavoriteByIDs() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("JrpRepository.RemoveFavoriteByIDs() : got =\n%v, want =\n%v", got, tt.want)
			}
			if err == nil && tt.args.ids != nil && len(tt.args.ids) > 0 {
				for _, id := range tt.args.ids {
					isExist, err := jrpChecker.IsExist(jrpDBFilePath, id)
					if err != nil {
						t.Errorf("JrpChecker.IsExist() : error =\n%v", err)
					}
					if isExist {
						isFavorited, err := jrpChecker.IsFavorited(jrpDBFilePath, id)
						if err != nil {
							t.Errorf("JrpChecker.IsFavorited() : error =\n%v", err)
						}
						if isFavorited {
							t.Errorf("JrpRepository.RemoveFavoriteByIDs() : is still favorited \n[%v]", id)
						}
					}
				}
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_RemoveFavoriteAll(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	jrpRepository := New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	timeProxy := timeproxy.New()

	type fields struct {
		FmtProxy      fmtproxy.Fmt
		SortProxy     sortproxy.Sort
		SqlProxy      sqlproxy.Sql
		StringsProxy  stringsproxy.Strings
		JrpDBFilePath string
	}
	type args struct {
		jrpDBFilePath string
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		want                  RemoveStatus
		wantFavoritedJrpCount int
		wantErr               bool
		setup                 func(*gomock.Controller, *fields)
		cleanup               func()
	}{
		{
			name: "positive testing (there is no jrps in the database)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedNone,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database and it is not favorited)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedNone,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is one jrp in the database and it is favorited)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedSuccessfully,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is two jrps in the database and both are not favorited)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedNone,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is two jrps in the database and both are favorited)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedSuccessfully,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (there is two jrps in the database and the one is favorited and the other is not favorited)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      sqlproxy.New(),
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedSuccessfully,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedFailed,
			wantFavoritedJrpCount: 1,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("SqlProxy.Open() failed"))
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (JrpRepository.createTableJrp() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedFailed,
			wantFavoritedJrpCount: 1,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Begin() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedFailed,
			wantFavoritedJrpCount: 1,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if _, err := jrpRepository.SaveHistory(jrpDBFilePath, []model.Jrp{
					{
						Phrase:    "test1",
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
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(nil, errors.New("DBInstance.Begin() failed"))
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (tx.Exec(query.RemoveAllFavorite) failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedFailed,
			wantFavoritedJrpCount: 1,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("TxInstance.Exec() failed"))
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (res.RowsAffected() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedFailed,
			wantFavoritedJrpCount: 1,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(0), errors.New("ResultInstance.RowsAffected() failed"))
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (tx.Commit() failed)",
			fields: fields{
				FmtProxy:      fmtproxy.New(),
				SortProxy:     sortproxy.New(),
				SqlProxy:      nil,
				StringsProxy:  stringsproxy.New(),
				JrpDBFilePath: jrpDBFilePath,
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:                  RemovedFailed,
			wantFavoritedJrpCount: 1,
			wantErr:               true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
					},
				); err != nil {
					t.Errorf("JrpRepository.SaveHistory() : error =\n%v", err)
				}
				if _, err := jrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{1}); err != nil {
					t.Errorf("JrpRepository.AddFavoriteByIDs() : error =\n%v", err)
				}
				mockResultInstance := mocksqlproxy.NewMockResultInstanceInterface(mockCtrl)
				mockResultInstance.EXPECT().RowsAffected().Return(int64(1), nil)
				mockTxInstance := mocksqlproxy.NewMockTxInstanceInterface(mockCtrl)
				mockTxInstance.EXPECT().Exec(gomock.Any()).Return(mockResultInstance, nil)
				mockTxInstance.EXPECT().Commit().Return(errors.New("TxInstance.Commit() failed"))
				mockTxInstance.EXPECT().Rollback().Return(nil)
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(&sqlproxy.ResultInstance{}, nil)
				mockDBInstance.EXPECT().Begin().Return(mockTxInstance, nil)
				mockDBInstance.EXPECT().Close().Return(nil)
				mockSqlProxy := mocksqlproxy.NewMockSql(mockCtrl)
				mockSqlProxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mockSqlProxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.RemoveFavoriteAll(tt.args.jrpDBFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.RemoveFavoriteAll() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("JrpRepository.RemoveFavoriteAll() : got =\n%v, want =\n%v", got, tt.want)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(tt.args.jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if len(favoritedJrps) != tt.wantFavoritedJrpCount {
				t.Errorf("JrpRepository.RemoveFavoriteAll() : len(favoritedJrps) =\n%v, wantFavoritedJrpCount =\n%v", favoritedJrps, tt.wantFavoritedJrpCount)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpRepository_createTableJrp(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathProxy,
		osProxy,
		userproxy.New(),
	)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, JRP_DB_FILE_NAME)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		db sqlproxy.DBInstanceInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(*gomock.Controller, *fields)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				db: nil,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (db.Exec() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				db: nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				mockDBInstance := mocksqlproxy.NewMockDBInstanceInterface(mockCtrl)
				mockDBInstance.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("DBInstance.Exec() failed"))
				mocksqlproxy := mocksqlproxy.NewMockSql(mockCtrl)
				mocksqlproxy.EXPECT().Open(gomock.Any(), gomock.Any()).Return(mockDBInstance, nil)
				tt.SqlProxy = mocksqlproxy
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
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
			j := New(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StringsProxy,
			)
			db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
			if err != nil {
				t.Errorf("SqlProxy.Open() : error =\n%v", err)
			}
			tt.args.db = db
			_, err = j.createTableJrp(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpRepository.createTableJrp() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
