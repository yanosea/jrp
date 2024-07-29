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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mu := mock_logic.NewMockUser(ctrl)
	mu.EXPECT().Current().Return(nil, errors.New("some error"))
	type args struct {
		e      Env
		u      User
		jrpEnv string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing (no env)",
			args: args{
				e:      OsEnv{},
				u:      tu,
				jrpEnv: "",
			},
			want:    filepath.Join(tcu.HomeDir, ".local", "share", "jrp"),
			wantErr: false,
		}, {
			name: "positive testing (with env)",
			args: args{
				e:      OsEnv{},
				u:      OsUser{},
				jrpEnv: filepath.Join(tcu.HomeDir, "jrp"),
			},
			want:    filepath.Join(tcu.HomeDir, "jrp"),
			wantErr: false,
		}, {
			name: "negative testing (user.Current() fails)",
			args: args{
				e:      OsEnv{},
				u:      mu,
				jrpEnv: filepath.Join(tcu.HomeDir, "jrp"),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.jrpEnv != "" {
				os.Setenv(constant.JRP_ENV_WORDNETJP_DIR, tt.args.jrpEnv)
				defer os.Unsetenv(constant.JRP_ENV_WORDNETJP_DIR)
			}
			got, err := GetDBFileDirPath(tt.args.e, tt.args.u)
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
