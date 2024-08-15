package util

import (
	"errors"
	"fmt"
	"os"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/internal/buffer"

	mock_buffer "github.com/yanosea/jrp/mock/buffer"
)

func TestCaptureOutput(t *testing.T) {
	type args struct {
		t        *testing.T
		fnc      func()
		capturer *DefaultCapturer
	}
	tests := []struct {
		name       string
		args       args
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name: "positive testing",
			args: args{
				t: t,
				fnc: func() {
					fmt.Println("stdout")
					fmt.Fprintln(os.Stderr, "stderr")
				},
				capturer: nil,
			},
			wantStdOut: "stdout\n",
			wantStdErr: "stderr\n",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				capturer := NewCapturer(&buffer.DefaultBuffer{}, &buffer.DefaultBuffer{})
				tt.capturer = capturer
			},
		}, {
			name: "negative testing (d.OutBuffer.ReadFrom(rOut) failed)",
			args: args{
				t: t,
				fnc: func() {
					fmt.Println("stdout")
					fmt.Fprintln(os.Stderr, "stderr")
				},
				capturer: nil,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mob := mock_buffer.NewMockBuffer(mockCtrl)
				mob.EXPECT().ReadFrom(gomock.Any()).Return(int64(0), errors.New("failed to read stdout buffer"))
				capturer := NewCapturer(mob, &buffer.DefaultBuffer{})
				tt.capturer = capturer
			},
		}, {
			name: "negative testing (d.ErrBuffer.ReadFrom(rErr) failed)",
			args: args{
				t: t,
				fnc: func() {
					fmt.Println("stdout")
					fmt.Fprintln(os.Stderr, "stderr")
				},
				capturer: nil,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				meb := mock_buffer.NewMockBuffer(mockCtrl)
				meb.EXPECT().ReadFrom(gomock.Any()).Return(int64(0), errors.New("failed to read stderr buffer"))
				capturer := NewCapturer(&buffer.DefaultBuffer{}, meb)
				tt.capturer = capturer
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

			gotStdOut, gotStdErr, err := tt.args.capturer.CaptureOutput(tt.args.t, tt.args.fnc)
			if (err != nil) != tt.wantErr {
				t.Errorf("CaptureOutput() : err = %v, wantErr %v", err, tt.wantErr)
			}
			if gotStdOut != "" && gotStdOut != tt.wantStdOut {
				t.Errorf("CaptureOutput() : gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != "" && gotStdErr != tt.wantStdErr {
				t.Errorf("CaptureOutput() : gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}
