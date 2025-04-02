package utility

import (
	"errors"
	"fmt"
	o "os"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	os := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()

	type args struct {
		os        proxy.Os
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
				os:        os,
				sdtBuffer: stdBuffer,
				errBuffer: errBuffer,
			},
			want: &capturer{
				os:        os,
				stdBuffer: stdBuffer,
				errBuffer: errBuffer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCapturer(tt.args.os, tt.args.sdtBuffer, tt.args.errBuffer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapturer_CaptureOutput(t *testing.T) {
	os := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()

	type fields struct {
		os        proxy.Os
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
				os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					if _, err := fmt.Fprint(o.Stdout, "stdout"); err != nil {
						t.Errorf("failed to write to stdout: %v", err)
					}
					if _, err := fmt.Fprint(o.Stderr, "stderr"); err != nil {
						t.Errorf("failed to write to stderr: %v", err)
					}
				},
			},
			wantStdOut: "stdout",
			wantStdErr: "stderr",
			wantErr:    false,
			setup:      nil,
		},
		{
			name: "negative testing (rOut, wOut, err := c.os.Pipe() failed)",
			fields: fields{
				os:        nil,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Pipe().Return(
					nil, nil, errors.New("Os.Pipe() failed"),
				)
				tt.os = mockOs
			},
		},
		{
			name: "negative testing (rErr, wErr, err := c.os.Pipe() failed)",
			fields: fields{
				os:        nil,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockFile1 := proxy.NewMockFile(mockCtrl)
				mockFile2 := proxy.NewMockFile(mockCtrl)
				mockOs.EXPECT().Pipe().Return(mockFile1, mockFile2, nil)
				mockOs.EXPECT().Pipe().Return(
					nil, nil, errors.New("Os.Pipe() failed"),
				)
				tt.os = mockOs
			},
		},
		{
			name: "negative testing (wOut.Close() failed)",
			fields: fields{
				os:        nil,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockStdoutReader := proxy.NewMockFile(mockCtrl)
				mockStdoutWriter := proxy.NewMockFile(mockCtrl)
				mockStdoutReader.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockStdoutWriter.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockStdoutWriter.EXPECT().Close().Return(errors.New("Close() failed"))
				mockOs.EXPECT().Pipe().Return(mockStdoutReader, mockStdoutWriter, nil)
				mockStderrReader := proxy.NewMockFile(mockCtrl)
				mockStderrWriter := proxy.NewMockFile(mockCtrl)
				mockStderrReader.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockStderrWriter.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockOs.EXPECT().Pipe().Return(mockStderrReader, mockStderrWriter, nil)
				tt.os = mockOs
			},
		},
		{
			name: "negative testing (wErr.Close() failed)",
			fields: fields{
				os:        nil,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockStdoutReader := proxy.NewMockFile(mockCtrl)
				mockStdoutWriter := proxy.NewMockFile(mockCtrl)
				mockStdoutReader.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockStdoutWriter.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockStdoutWriter.EXPECT().Close().Return(nil)
				mockOs.EXPECT().Pipe().Return(mockStdoutReader, mockStdoutWriter, nil)
				mockStderrReader := proxy.NewMockFile(mockCtrl)
				mockStderrWriter := proxy.NewMockFile(mockCtrl)
				mockStderrReader.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockStderrWriter.EXPECT().AsOsFile().Return(nil).AnyTimes()
				mockStderrWriter.EXPECT().Close().Return(errors.New("Close() failed"))
				mockOs.EXPECT().Pipe().Return(mockStderrReader, mockStderrWriter, nil)
				tt.os = mockOs
			},
		},
		{
			name: "negative testing (c.OutBuffer.ReadFrom(rOut) failed)",
			fields: fields{
				os:        os,
				StdBuffer: nil,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					if _, err := fmt.Fprint(o.Stdout, "stdout"); err != nil {
						t.Errorf("failed to write to stdout: %v", err)
					}
					if _, err := fmt.Fprint(o.Stderr, "stderr"); err != nil {
						t.Errorf("failed to write to stderr: %v", err)
					}
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
				os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: nil,
			},
			args: args{
				fnc: func() {
					if _, err := fmt.Fprint(o.Stdout, "stdout"); err != nil {
						t.Errorf("failed to write to stdout: %v", err)
					}
					if _, err := fmt.Fprint(o.Stderr, "stderr"); err != nil {
						t.Errorf("failed to write to stderr: %v", err)
					}
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
				os:        tt.fields.os,
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
