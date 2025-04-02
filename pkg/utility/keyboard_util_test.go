package utility

import (
	"errors"
	"reflect"
	"testing"

	"github.com/eiannone/keyboard"

	"github.com/yanosea/jrp/v2/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewKeyboardUtil(t *testing.T) {
	keyboard := proxy.NewKeyboard()

	type args struct {
		keyboard proxy.Keyboard
	}
	tests := []struct {
		name string
		args args
		want KeyboardUtil
	}{
		{
			name: "positive testing",
			args: args{
				keyboard: keyboard,
			},
			want: &keyboardUtil{
				keyboard: keyboard,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKeyboardUtil(tt.args.keyboard); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeyboardUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keyboardUtil_CloseKeyboard(t *testing.T) {
	type fields struct {
		keyboard proxy.Keyboard
	}
	tests := []struct {
		name   string
		fields fields
		setup  func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				keyboard: nil,
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().Close().Return()
				tt.keyboard = mockKeyboard
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			k := &keyboardUtil{
				keyboard: tt.fields.keyboard,
			}
			if err := k.CloseKeyboard(); err != nil {
				t.Errorf("keyboardUtil.CloseKeyboard() error = %v", err)
			}
		})
	}
}

func Test_keyboardUtil_GetKey(t *testing.T) {
	type fields struct {
		keyboard proxy.Keyboard
	}
	type args struct {
		timeoutSec int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				keyboard: nil,
			},
			args: args{
				timeoutSec: 1,
			},
			want:    "a",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().GetKey(1).Return('a', keyboard.KeyEnter, nil)
				tt.keyboard = mockKeyboard
			},
		},
		{
			name: "negative testing (k.keyboard.GetKey() failed)",
			fields: fields{
				keyboard: nil,
			},
			args: args{
				timeoutSec: 1,
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().GetKey(1).Return(rune(0), keyboard.Key(0), errors.New("KeyboardProxy.GetKey() failed"))
				tt.keyboard = mockKeyboard
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			k := &keyboardUtil{
				keyboard: tt.fields.keyboard,
			}
			got, err := k.GetKey(tt.args.timeoutSec)
			if (err != nil) != tt.wantErr {
				t.Errorf("keyboardUtil.GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("keyboardUtil.GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keyboardUtil_OpenKeyboard(t *testing.T) {
	type fields struct {
		keyboard proxy.Keyboard
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				keyboard: nil,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().Open().Return(nil)
				tt.keyboard = mockKeyboard
			},
		},
		{
			name: "negative testing (k.keyboard.Open() failed)",
			fields: fields{
				keyboard: nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockKeyboard := proxy.NewMockKeyboard(mockCtrl)
				mockKeyboard.EXPECT().Open().Return(errors.New("KeyboardProxy.Open() failed"))
				tt.keyboard = mockKeyboard
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			k := &keyboardUtil{
				keyboard: tt.fields.keyboard,
			}
			if err := k.OpenKeyboard(); (err != nil) != tt.wantErr {
				t.Errorf("keyboardUtil.OpenKeyboard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
