package keyboardproxy

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New()
		})
	}
}

func TestKeyboardProxy_Open(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "positive testing",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyboardProxy := New()
			if err := keyboardProxy.Open(); (err != nil) != tt.wantErr {
				t.Errorf("KeyboardProxy.Open() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
		})
	}
}

func TestKeyboardProxy_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyboardProxy := New()
			if err := keyboardProxy.Open(); err != nil {
				t.Errorf("KeyboardProxy.Open() : error =\n%v", err)
				return
			}
			keyboardProxy.Close()
		})
	}
}

func TestKeyboardProxy_GetKey(t *testing.T) {
	keyboardProxy := New()
	type arg struct {
		timeoutSec int
	}
	tests := []struct {
		name    string
		args    arg
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name:    "positive testing",
			args:    arg{timeoutSec: 0},
			wantErr: false,
			setup: func() {
				if err := keyboardProxy.Open(); err != nil {
					t.Errorf("KeyboardProxy.Open() : error =\n%v", err)
				}
			},
			cleanup: func() {
				keyboardProxy.Close()
			},
		},
		{
			name:    "negative testing",
			args:    arg{timeoutSec: 10},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			_, _, err := keyboardProxy.GetKey(tt.args.timeoutSec)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyboardProxy.GetKey() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
