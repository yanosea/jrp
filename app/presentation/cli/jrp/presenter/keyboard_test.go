package presenter

import (
	"testing"

	"github.com/eiannone/keyboard"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestCloseKeyboard(t *testing.T) {
	origKu := Ku

	tests := []struct {
		name    string
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing",
			setup: func(mockCtrl *gomock.Controller) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().Close()
				Ku = utility.NewKeyboardUtil(mockKeyboard)
			},
			cleanup: func() {
				Ku = origKu
			},
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		if tt.setup != nil {
			tt.setup(mockCtrl)
		}
		defer func() {
			if tt.cleanup != nil {
				tt.cleanup()
			}
		}()
		t.Run(tt.name, func(t *testing.T) {
			if err := CloseKeyboard(); err != nil {
				t.Errorf("CloseKeyboard() error = %v", err)
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	origKu := Ku

	type args struct {
		timeoutSec int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				timeoutSec: 1,
			},
			want:    "a",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().GetKey(1).Return('a', keyboard.Key(0), nil)
				Ku = utility.NewKeyboardUtil(mockKeyboard)
			},
			cleanup: func() {
				Ku = origKu
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			got, err := GetKey(tt.args.timeoutSec)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenKeyboard(t *testing.T) {
	orgK := Ku

	tests := []struct {
		name    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name:    "positive testing",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().Open().Return(nil)
				Ku = utility.NewKeyboardUtil(mockKeyboard)
			},
			cleanup: func() {
				Ku = orgK
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := OpenKeyboard(); (err != nil) != tt.wantErr {
				t.Errorf("OpenKeyboard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
