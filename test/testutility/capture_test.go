package testutility

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strconv"

	"github.com/yanosea/jrp/mock/app/proxy/buffer"
	"go.uber.org/mock/gomock"
)

func TestNewCapturer(t *testing.T) {
	outBufferProxy := bufferproxy.New()
	errorBufferProxy := bufferproxy.New()
	osProxy := osproxy.New()

	type args struct {
		outBuffer   bufferproxy.Buffer
		errorBuffer bufferproxy.Buffer
		osProxy     osproxy.Os
	}
	tests := []struct {
		name string
		args args
		want *Capturer
	}{
		{
			name: "positive testing",
			args: args{
				outBuffer:   outBufferProxy,
				errorBuffer: errorBufferProxy,
				osProxy:     osProxy,
			},
			want: &Capturer{
				OutBuffer: outBufferProxy,
				ErrBuffer: errorBufferProxy,
				OsProxy:   osProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCapturer(
				tt.args.outBuffer,
				tt.args.errorBuffer,
				tt.args.osProxy,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCapturer() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestCapturer_CaptureOutput(t *testing.T) {
	util := utility.New(
		fmtproxy.New(),
		osproxy.New(),
		strconvproxy.New(),
	)

	type fields struct {
		outBuffer   bufferproxy.Buffer
		errorBuffer bufferproxy.Buffer
		osProxy     osproxy.Os
	}
	type args struct {
		t   *testing.T
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
				outBuffer:   bufferproxy.New(),
				errorBuffer: bufferproxy.New(),
				osProxy:     osproxy.New(),
			},
			args: args{
				t: t,
				fnc: func() {
					util.PrintlnWithWriter(osproxy.Stdout, "stdout")
					util.PrintlnWithWriter(osproxy.Stderr, "stderr")
				},
			},
			wantStdOut: "stdout\n",
			wantStdErr: "stderr\n",
			wantErr:    false,
			setup:      nil,
		},
		{
			name: "negative testing (c.OutBuffer.ReadFrom(rOut) failed)",
			fields: fields{
				outBuffer:   nil,
				errorBuffer: bufferproxy.New(),
				osProxy:     osproxy.New(),
			},
			args: args{
				t: t,
				fnc: func() {
					util.PrintlnWithWriter(osproxy.Stdout, "stdout")
					util.PrintlnWithWriter(osproxy.Stderr, "stderr")
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockBufferProxy := mockbufferproxy.NewMockBuffer(mockCtrl)
				mockBufferProxy.EXPECT().ReadFrom(gomock.Any()).Return(
					int64(0),
					errors.New("BufferProxy.ReadFrom() failed"),
				)
				tt.outBuffer = mockBufferProxy
			},
		},
		{
			name: "negative testing (c.ErrBuffer.ReadFrom(rErr) failed)",
			fields: fields{
				outBuffer:   bufferproxy.New(),
				errorBuffer: nil,
				osProxy:     osproxy.New(),
			},
			args: args{
				t: t,
				fnc: func() {
					util.PrintlnWithWriter(osproxy.Stdout, "stdout")
					util.PrintlnWithWriter(osproxy.Stderr, "stderr")
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockBufferProxy := mockbufferproxy.NewMockBuffer(mockCtrl)
				mockBufferProxy.EXPECT().ReadFrom(gomock.Any()).Return(
					int64(0),
					errors.New("BufferProxy.ReadFrom() failed"),
				)
				tt.errorBuffer = mockBufferProxy
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
			c := NewCapturer(tt.fields.outBuffer, tt.fields.errorBuffer, tt.fields.osProxy)
			gotStdOut, gotStdErr, err := c.CaptureOutput(tt.args.t, tt.args.fnc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Capturer.CaptureOutput() : gotStdOut =\n%v, want =\n%v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Capturer.CaptureOutput() : gotStdErr =\n%v, want =\n%v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}
