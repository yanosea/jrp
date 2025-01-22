package utility

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()

	type args struct {
		sdtBuffer proxy.Buffer
		errBuffer proxy.Buffer
	}
	tests := []struct {
		name string
		args args
		want *capturer
	}{
		{
			name: "positive testing",
			args: args{
				sdtBuffer: stdBuffer,
				errBuffer: errBuffer,
			},
			want: &capturer{
				stdBuffer: stdBuffer,
				errBuffer: errBuffer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCapturer(tt.args.sdtBuffer, tt.args.errBuffer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapturer_CaptureOutput(t *testing.T) {
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()

	type fields struct {
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func()
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					fmt.Fprint(os.Stdout, "stdout")
					fmt.Fprint(os.Stderr, "stderr")
				},
			},
			wantStdOut: "stdout",
			wantStdErr: "stderr",
			wantErr:    false,
			setup:      nil,
		},
		{
			name: "negative testing (c.OutBuffer.ReadFrom(rOut) failed)",
			fields: fields{
				StdBuffer: nil,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					fmt.Fprint(os.Stdout, "stdout")
					fmt.Fprint(os.Stderr, "stderr")
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockBuffer := proxy.NewMockBuffer(mockCtrl)
				mockBuffer.EXPECT().ReadFrom(gomock.Any()).Return(
					int64(0),
					errors.New("BufferProxy.ReadFrom() failed"),
				)
				tt.StdBuffer = mockBuffer
			},
		},
		{
			name: "negative testing (c.ErrBuffer.ReadFrom(rErr) failed)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: nil,
			},
			args: args{
				fnc: func() {
					fmt.Fprint(os.Stdout, "stdout")
					fmt.Fprint(os.Stderr, "stderr")
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockBuffer := proxy.NewMockBuffer(mockCtrl)
				mockBuffer.EXPECT().ReadFrom(gomock.Any()).Return(
					int64(0),
					errors.New("BufferProxy.ReadFrom() failed"),
				)
				tt.ErrBuffer = mockBuffer
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
			c := &capturer{
				stdBuffer: tt.fields.StdBuffer,
				errBuffer: tt.fields.ErrBuffer,
			}
			gotStdOut, gotStdErr, err := c.CaptureOutput(tt.args.fnc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Capturer.CaptureOutput() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Capturer.CaptureOutput() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}
