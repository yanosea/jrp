package testutility

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/database/jrp/repository"
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

	"github.com/yanosea/jrp/mock/app/proxy/sql"
	"go.uber.org/mock/gomock"
)

func TestNewJrpChecker(t *testing.T) {
	fmtProxy := fmtproxy.New()
	sortProxy := sortproxy.New()
	sqlProxy := sqlproxy.New()
	strconvProxy := strconvproxy.New()
	stringProxy := stringsproxy.New()

	type args struct {
		fmtProxy     fmtproxy.Fmt
		sortProxy    sortproxy.Sort
		sqlProxy     sqlproxy.Sql
		strconvProxy strconvproxy.Strconv
		stringsProxy stringsproxy.Strings
	}
	tests := []struct {
		name string
		args args
		want *JrpChecker
	}{
		{
			name: "positive testing",
			args: args{
				fmtProxy:     fmtProxy,
				sortProxy:    sortProxy,
				sqlProxy:     sqlProxy,
				strconvProxy: strconvProxy,
				stringsProxy: stringProxy,
			},
			want: &JrpChecker{
				FmtProxy:     fmtProxy,
				SortProxy:    sortProxy,
				SqlProxy:     sqlProxy,
				StrconvProxy: strconvProxy,
				StringsProxy: stringProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJrpChecker(tt.args.fmtProxy, tt.args.sortProxy, tt.args.sqlProxy, tt.args.strconvProxy, tt.args.stringsProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJrpChecker() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestJrpChecker_GetJrpSeq(t *testing.T) {
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
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := repository.New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StrconvProxy strconvproxy.Strconv
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
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
				StrconvProxy: strconvproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    0,
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
				if _, err := jrpRepository.RemoveHistoryAll(jrpDBFilePath, true); err != nil {
					t.Errorf("JrpRepository.RemoveHistoryAll() : error =\n%v", err)
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
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    1,
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
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    0,
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
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    0,
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
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
			},
			want:    0,
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
			j := NewJrpChecker(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StrconvProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.GetJrpSeq(tt.args.jrpDBFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpChecker.GetJrpSeq() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JrpChecker.GetJrpSeq() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpChecker_IsExist(t *testing.T) {
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
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := repository.New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StrconvProxy strconvproxy.Strconv
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		id            int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (the jrp does not exist)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
				StrconvProxy: strconvproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            2,
			},
			want:    false,
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
			name: "positive testing (the jrp exists)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
				StrconvProxy: strconvproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    true,
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
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    false,
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
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    false,
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
				mockDBInstance.EXPECT().Query(gomock.Any(), "1").Return(nil, errors.New("DBInstance.Query() failed"))
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
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    false,
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
				mockDBInstance.EXPECT().Query(gomock.Any(), "1").Return(mockRowsInstance, nil)
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
			j := NewJrpChecker(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StrconvProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.IsExist(tt.args.jrpDBFilePath, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpChecker.IsExist() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JrpChecker.IsExist() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpChecker_IsFavorited(t *testing.T) {
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
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()
	jrpRepository := repository.New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StrconvProxy strconvproxy.Strconv
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		jrpDBFilePath string
		id            int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (the jrp does not exist)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
				StrconvProxy: strconvproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            2,
			},
			want:    false,
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
			name: "positive testing (the jrp exists)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StringsProxy: stringsproxy.New(),
				StrconvProxy: strconvproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    false,
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
			name: "negative testing (SqlProxy.Open() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    false,
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
			name: "negative testing (db.Query() failed)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     nil,
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    false,
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
				mockDBInstance.EXPECT().Query(gomock.Any(), "1").Return(nil, errors.New("DBInstance.Query() failed"))
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
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				jrpDBFilePath: jrpDBFilePath,
				id:            1,
			},
			want:    false,
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
				mockDBInstance.EXPECT().Query(gomock.Any(), "1").Return(mockRowsInstance, nil)
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
			j := NewJrpChecker(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StrconvProxy,
				tt.fields.StringsProxy,
			)
			got, err := j.IsFavorited(tt.args.jrpDBFilePath, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("JrpChecker.IsFavorited() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JrpChecker.IsFavorited() : got =\n%v, want =\n%v", got, tt.want)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestJrpChecker_IsSameJrps(t *testing.T) {
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()

	type fields struct {
		FmtProxy     fmtproxy.Fmt
		SortProxy    sortproxy.Sort
		SqlProxy     sqlproxy.Sql
		StrconvProxy strconvproxy.Strconv
		StringsProxy stringsproxy.Strings
	}
	type args struct {
		got  []model.Jrp
		want []model.Jrp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "positive testing(len is not same)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				got: []model.Jrp{
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
				want: []model.Jrp{},
			},
			want: false,
		},
		{
			name: "positive testing(favorited and not same)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				got: []model.Jrp{
					{
						ID:          1,
						Phrase:      "test1",
						Prefix:      sqlProxy.StringToNullString("prefix"),
						Suffix:      sqlProxy.StringToNullString(""),
						IsFavorited: 1,
						CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
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
			},
			want: false,
		},
		{
			name: "positive testing(favorited and same)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				got: []model.Jrp{
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
				want: []model.Jrp{
					{
						ID:          1,
						Phrase:      "test1",
						Prefix:      sqlProxy.StringToNullString(""),
						Suffix:      sqlProxy.StringToNullString(""),
						IsFavorited: 1,
						CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt:   timeProxy.Date(0001, 01, 01, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
			},
			want: true,
		},
		{
			name: "positive testing(not favorited and not same)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				got: []model.Jrp{
					{
						ID:          1,
						Phrase:      "test1",
						Prefix:      sqlProxy.StringToNullString("prefix"),
						Suffix:      sqlProxy.StringToNullString(""),
						IsFavorited: 0,
						CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
					},
				},
				want: []model.Jrp{
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
			},
			want: false,
		},
		{
			name: "positive testing(not favorited and same)",
			fields: fields{
				FmtProxy:     fmtproxy.New(),
				SortProxy:    sortproxy.New(),
				SqlProxy:     sqlproxy.New(),
				StrconvProxy: strconvproxy.New(),
				StringsProxy: stringsproxy.New(),
			},
			args: args{
				got: []model.Jrp{
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
				want: []model.Jrp{
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
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJrpChecker(
				tt.fields.FmtProxy,
				tt.fields.SortProxy,
				tt.fields.SqlProxy,
				tt.fields.StrconvProxy,
				tt.fields.StringsProxy,
			)
			if got := j.IsSameJrps(tt.args.got, tt.args.want); got != tt.want {
				t.Errorf("JrpChecker.IsSameJrps() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func Test_isSameNullStringInstance(t *testing.T) {
	type args struct {
		a *sqlproxy.NullStringInstance
		b *sqlproxy.NullStringInstance
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive testing (both are nil)",
			args: args{
				a: nil,
				b: nil,
			},
			want: true,
		},
		{
			name: "positive testing (one is nil, the other is not nil)",
			args: args{
				a: nil,
				b: &sqlproxy.NullStringInstance{},
			},
			want: false,
		},
		{
			name: "positive testing (not the same string)",
			args: args{
				a: &sqlproxy.NullStringInstance{
					FieldNullString: &sql.NullString{
						String: "test1",
						Valid:  true,
					},
				},
				b: &sqlproxy.NullStringInstance{
					FieldNullString: &sql.NullString{
						String: "test2",
						Valid:  true,
					},
				},
			},
			want: false,
		},
		{
			name: "positive testing (the same string)",
			args: args{
				a: &sqlproxy.NullStringInstance{
					FieldNullString: &sql.NullString{
						String: "test1",
						Valid:  true,
					},
				},
				b: &sqlproxy.NullStringInstance{
					FieldNullString: &sql.NullString{
						String: "test1",
						Valid:  true,
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSameNullStringInstance(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("isSameNullStringInstance() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func Test_isSameTimeInstance(t *testing.T) {
	timeProxy := timeproxy.New()

	type args struct {
		a *timeproxy.TimeInstance
		b *timeproxy.TimeInstance
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive testing (both are nil)",
			args: args{
				a: nil,
				b: nil,
			},
			want: true,
		},
		{
			name: "positive testing (one is nil, the other is not nil)",
			args: args{
				a: nil,
				b: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
			},
			want: false,
		},
		{
			name: "positive testing (not the same time)",
			args: args{
				a: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				b: timeProxy.Date(0001, 01, 01, 0, 0, 0, 0, &timeproxy.UTC),
			},
			want: false,
		},
		{
			name: "positive testing (the same time)",
			args: args{
				a: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
				b: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSameTimeInstance(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("isSameTimeInstance() : =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}
