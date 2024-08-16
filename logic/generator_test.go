package logic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/database"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/spinnerservice"
	"github.com/yanosea/jrp/internal/usermanager"

	mock_database "github.com/yanosea/jrp/mock/database"
	mock_fs "github.com/yanosea/jrp/mock/fs"
	mock_usermanager "github.com/yanosea/jrp/mock/usermanager"
)

func TestNewJapaneseRandomPhraseGenerator(t *testing.T) {
	type args struct {
		u usermanager.UserProvider
		d database.DatabaseProvider
		f fs.FileManager
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{u: usermanager.OSUserProvider{}, d: database.SQLiteProvider{}, f: fs.OsFileManager{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewJapaneseRandomPhraseGenerator(tt.args.u, tt.args.d, tt.args.f)
			if u == nil {
				t.Errorf("NewJapaneseRandomPhraseGenerator() : returned nil")
			}
		})
	}
}

func TestDefineNumber(t *testing.T) {
	type args struct {
		num    int
		argNum string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "positive testing (num is -1, argNum is empty)",
			args: args{num: -1, argNum: ""},
			want: 1,
		},
		{
			name: "positive testing (num is 0, argNum is empty)",
			args: args{num: 0, argNum: ""},
			want: 1,
		}, {
			name: "positive testing (num is 1, argNum is empty)",
			args: args{num: 1, argNum: ""},
			want: 1,
		}, {
			name: "positive testing (num is 2, argsNum is empty)",
			args: args{num: 2, argNum: ""},
			want: 2,
		}, {
			name: "positive testing (num is 0, argNum is -1)",
			args: args{num: 0, argNum: "-1"},
			want: 1,
		}, {
			name: "positive testing (num is 0, argNum is 0)",
			args: args{num: 0, argNum: "0"},
			want: 1,
		}, {
			name: "positive testing (num is 0, argNum can't be converted to int)",
			args: args{num: 0, argNum: "test"},
			want: 1,
		}, {
			name: "positive testing (num is 0, argNum is 1)",
			args: args{num: 0, argNum: "1"},
			want: 1,
		}, {
			name: "positive testing (num is 0, argNum is 2)",
			args: args{num: 0, argNum: "2"},
			want: 2,
		}, {
			name: "positive testing (num is 3, argNum is 2)",
			args: args{num: 3, argNum: "2"},
			want: 3,
		}, {
			name: "positive testing (num is 2, argNum is 3)",
			args: args{num: 2, argNum: "3"},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefineNumber(tt.args.num, tt.args.argNum); got != tt.want {
				t.Errorf("DefineNumber() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tu := usermanager.OSUserProvider{}
	tcu, _ := tu.Current()
	defaultDBFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	tdl := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{}, spinnerservice.NewRealSpinnerService())
	if err := tdl.Download(); err != nil {
		t.Error(err)
	}

	type args struct {
		generator *JapaneseRandomPhraseGenerator
		num       int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name:    "positive testing (num is 1)",
			args:    args{generator: nil, num: 1},
			want:    1,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})
				tt.generator = generator
			},
		}, {
			name:    "positive testing (num is 2)",
			args:    args{generator: nil, num: 2},
			want:    2,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})
				tt.generator = generator
			},
		}, {
			name:    "positive testing (DBFile does not exist)",
			args:    args{generator: nil, num: 1},
			want:    0,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mfs := mock_fs.NewMockFileManager(mockCtrl)
				mfs.EXPECT().Exists(filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME)).Return(false)
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, mfs)
				tt.generator = generator
			},
		}, {
			name:    "negative testing (GetFileDirPath() fails)",
			args:    args{generator: nil, num: 1},
			want:    0,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mu := mock_usermanager.NewMockUserProvider(mockCtrl)
				mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))
				generator := NewJapaneseRandomPhraseGenerator(mu, database.SQLiteProvider{}, fs.OsFileManager{})
				tt.generator = generator
			},
		}, {
			name:    "negative testing (Connect() fails)",
			args:    args{generator: nil, num: 1},
			want:    0,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mdb := mock_database.NewMockDatabaseProvider(mockCtrl)
				mdb.EXPECT().Connect(gomock.Any()).Return(nil, errors.New("failed to connect db"))
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, mdb, fs.OsFileManager{})
				tt.generator = generator
			},
		}, {
			name:    "negative testing (Query() fails)",
			args:    args{generator: nil, num: 1},
			want:    0,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mdb := mock_database.NewMockDatabaseProvider(mockCtrl)
				jrpg := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})
				db, _ := jrpg.DbProvider.Connect(filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME))
				gomock.InOrder(
					mdb.EXPECT().Connect(filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME)).Return(db, nil),
					mdb.EXPECT().Query(db, constant.GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS).Return(nil, errors.New("failed to execute query")),
				)
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, mdb, fs.OsFileManager{})
				tt.generator = generator
			},
		}, {
			name:    "negative testing (rows.Scan(&word.Lemma, &word.Pos) fails)",
			args:    args{generator: nil, num: 1},
			want:    0,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mr := mock_database.NewMockRows(mockCtrl)
				gomock.InOrder(
					mr.EXPECT().Next().Return(true),
					mr.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(errors.New("failed to scan")),
					mr.EXPECT().Close().Return(nil),
				)
				var mrAsRows database.Rows = mr
				mdb := mock_database.NewMockDatabaseProvider(mockCtrl)
				jrpg := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})
				db, _ := jrpg.DbProvider.Connect(filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME))
				gomock.InOrder(
					mdb.EXPECT().Connect(filepath.Join(defaultDBFileDirPath, constant.WNJPN_DB_FILE_NAME)).Return(db, nil),
					mdb.EXPECT().Query(db, constant.GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS).Return(mrAsRows, nil),
				)
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, mdb, fs.OsFileManager{})
				tt.generator = generator
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.setup != nil {
				tt.setup(ctrl, &tt.args)
			}

			jrps, err := tt.args.generator.Generate(tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(jrps) != tt.want {
				t.Errorf("Generate() got(count) = %v, want = %v", len(jrps), tt.want)
			}
			for _, jrp := range jrps {
				fmt.Println(jrp)
			}
		})
	}
	os.RemoveAll(defaultDBFileDirPath)
}
