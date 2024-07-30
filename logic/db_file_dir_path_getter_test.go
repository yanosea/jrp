package logic

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/constant"
	mock_logic "github.com/yanosea/jrp/mock/logic"
)

func TestGetDBFileDirPath(t *testing.T) {
	tu := OsUser{}
	tcu, _ := tu.Current()

	type args struct {
		dbFileDirPathGetter *DBFileDirPathGetter
		env                 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name:    "positive testing (no env)",
			args:    args{dbFileDirPathGetter: nil, env: ""},
			want:    filepath.Join(tcu.HomeDir, ".local", "share", "jrp"),
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				dbFileDirPathGetter := NewDBFileDirPathGetter(OsEnv{}, OsUser{})
				tt.dbFileDirPathGetter = dbFileDirPathGetter
			},
		}, {
			name:    "positive testing (with env)",
			args:    args{dbFileDirPathGetter: nil, env: filepath.Join(tcu.HomeDir, "jrp")},
			want:    filepath.Join(tcu.HomeDir, "jrp"),
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				dbFileDirPathGetter := NewDBFileDirPathGetter(OsEnv{}, OsUser{})
				tt.dbFileDirPathGetter = dbFileDirPathGetter
			},
		}, {
			name:    "negative testing (user.Current() fails)",
			args:    args{dbFileDirPathGetter: nil, env: ""},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mu := mock_logic.NewMockUser(mockCtrl)
				mu.EXPECT().Current().Return(nil, errors.New("failed to get current user"))
				dbFileDirPathGetter := NewDBFileDirPathGetter(OsEnv{}, mu)
				tt.dbFileDirPathGetter = dbFileDirPathGetter
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

			if tt.args.env != "" {
				os.Setenv(constant.JRP_ENV_WORDNETJP_DIR, tt.args.env)
				defer os.Unsetenv(constant.JRP_ENV_WORDNETJP_DIR)
			}
			got, err := tt.args.dbFileDirPathGetter.GetFileDirPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDBFileDirPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDBFileDirPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
