package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/internal/db"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/rand"
	"github.com/yanosea/jrp/internal/usermanager"
)

func TestDefineNumber(t *testing.T) {
	type args struct {
		generator *JapaneseRandomPhraseGenerator
		num       int
		argNum    string
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
	tdl := NewDBFileDownloader(usermanager.OSUserProvider{}, fs.OsFileManager{}, httpclient.DefaultHTTPClient{}, iomanager.DefaultIOHelper{}, gzip.DefaultGzipHandler{})
	tdl.Download()

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
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, db.SQLiteProvider{}, fs.OsFileManager{}, rand.NewDefaultRandomGenerator())
				tt.generator = generator
			},
		}, {
			name:    "positive testing (num is 2)",
			args:    args{generator: nil, num: 2},
			want:    2,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				generator := NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, db.SQLiteProvider{}, fs.OsFileManager{}, rand.NewDefaultRandomGenerator())
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
	defaultDBFileDirPath := filepath.Join(tcu.HomeDir, ".local", "share", "jrp")
	os.RemoveAll(defaultDBFileDirPath)
}
