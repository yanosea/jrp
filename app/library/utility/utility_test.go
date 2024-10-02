package utility

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strconv"

	"github.com/yanosea/jrp/test/testutility"
)

func TestNew(t *testing.T) {
	fmtProxy := fmtproxy.New()
	osProxy := osproxy.New()
	strconvProxy := strconvproxy.New()

	type args struct {
		fmtProxy     fmtproxy.Fmt
		osProxy      osproxy.Os
		strconvProxy strconvproxy.Strconv
	}
	tests := []struct {
		name string
		args args
		want *Utility
	}{
		{
			name: "positive testing",
			args: args{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			want: &Utility{
				FmtProxy:     fmtProxy,
				OsProxy:      osProxy,
				StrconvProxy: strconvProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.fmtProxy, tt.args.osProxy, tt.args.strconvProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestUtility_PrintlnWithWriter(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
	util := New(
		fmtproxy.New(),
		osproxy.New(),
		strconvproxy.New(),
	)

	type fields struct {
		t        *testing.T
		fnc      func()
		capturer *testutility.Capturer
	}
	tests := []struct {
		name       string
		fields     fields
		wantStdOut string
		wantStdErr string
		wantErr    bool
	}{
		{
			name: "positive testing (stdout)",
			fields: fields{
				t: t,
				fnc: func() {
					util.PrintlnWithWriter(osproxy.Stdout, "stdout")
				},
				capturer: capturer,
			},
			wantStdOut: "stdout\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (stderr)",
			fields: fields{
				t: t,
				fnc: func() {
					util.PrintlnWithWriter(osproxy.Stderr, "stderr")
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "stderr\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdOut, gotStdErr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Utility.PrintlnWithWriter() : gotStdOut =\n%v, want =\n%v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Utility.PrintlnWithWriter() : gotStdErr =\n%v, want =\n%v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}

func TestUtility_GetMaxConvertibleString(t *testing.T) {
	fmtProxy := fmtproxy.New()
	osProxy := osproxy.New()
	strconvProxy := strconvproxy.New()

	type fields struct {
		fmtProxy     fmtproxy.Fmt
		osProxy      osproxy.Os
		strconvProxy strconvproxy.Strconv
	}
	type args struct {
		args []string
		def  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "positive testing (args have convertible string, it is first)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				args: []string{"2", "test", "3"},
				def:  "1",
			},
			want: "3",
		},
		{
			name: "positive testing (args have convertible string, it is second)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				args: []string{"test", "4", "5"},
				def:  "1",
			},
			want: "5",
		},
		{
			name: "positive testing (args have convertible string, it is last)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				args: []string{"test", "test", "6"},
				def:  "1",
			},
			want: "6",
		},
		{
			name: "positive testing (args have no convertible string)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				args: []string{"test", "test", "test"},
				def:  "1",
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(
				tt.fields.fmtProxy,
				tt.fields.osProxy,
				tt.fields.strconvProxy,
			)
			if got := u.GetMaxConvertibleString(tt.args.args, tt.args.def); got != tt.want {
				t.Errorf("Utility.GetMaxConvertibleString() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestUtility_GetLargerNumber(t *testing.T) {
	fmtProxy := fmtproxy.New()
	osProxy := osproxy.New()
	strconvProxy := strconvproxy.New()

	type fields struct {
		fmtProxy     fmtproxy.Fmt
		osProxy      osproxy.Os
		strconvProxy strconvproxy.Strconv
	}
	type args struct {
		num    int
		argNum string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "positive testing (num is -1, argNum is empty)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    -1,
				argNum: "",
			},
			want: 1,
		},
		{
			name: "positive testing (num is 0, argNum is empty)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    0,
				argNum: "",
			},
			want: 1,
		},
		{
			name: "positive testing (num is 1, argNum is empty)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    1,
				argNum: "",
			},
			want: 1,
		},
		{
			name: "positive testing (num is 2, argsNum is empty)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    2,
				argNum: "",
			},
			want: 2,
		},
		{
			name: "positive testing (num is 0, argNum is -1)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    0,
				argNum: "-1",
			},
			want: 1,
		},
		{
			name: "positive testing (num is 0, argNum is 0)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    0,
				argNum: "0",
			},
			want: 1,
		},
		{
			name: "positive testing (num is 0, argNum can't be converted to int)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    0,
				argNum: "test",
			},
			want: 1,
		},
		{
			name: "positive testing (num is 0, argNum is 1)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    0,
				argNum: "1",
			},
			want: 1,
		},
		{
			name: "positive testing (num is 0, argNum is 2)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    0,
				argNum: "2",
			},
			want: 2,
		},
		{
			name: "positive testing (num is 3, argNum is 2)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    3,
				argNum: "2",
			},
			want: 3,
		},
		{
			name: "positive testing (num is 2, argNum is 3)",
			fields: fields{
				fmtProxy:     fmtProxy,
				osProxy:      osProxy,
				strconvProxy: strconvProxy,
			},
			args: args{
				num:    2,
				argNum: "3",
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(
				tt.fields.fmtProxy,
				tt.fields.osProxy,
				tt.fields.strconvProxy,
			)
			if got := u.GetLargerNumber(tt.args.num, tt.args.argNum); got != tt.want {
				t.Errorf("Utility.GetLargerNumber() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestUtility_CreateDirIfNotExist(t *testing.T) {
	filepathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	testDirPath := filepathProxy.Join(osProxy.TempDir(), "jrp_test")

	type fields struct {
		fmtProxy     fmtproxy.Fmt
		osProxy      osproxy.Os
		strconvProxy strconvproxy.Strconv
	}
	type args struct {
		dirPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (dirPath does not exists)",
			fields: fields{
				fmtProxy:     fmtproxy.New(),
				osProxy:      osproxy.New(),
				strconvProxy: strconvproxy.New(),
			},
			args: args{
				dirPath: testDirPath,
			},
			wantErr: false,
			setup: func() {
				if err := osProxy.RemoveAll(testDirPath); err != nil {
					t.Errorf("OsProxy.RemoveAll(%v) : error =\n%v", testDirPath, err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(testDirPath); err != nil {
					t.Errorf("OsProxy.RemoveAll(%v) : error =\n%v", testDirPath, err)
				}
			},
		},
		{
			name: "positive testing (dirPath exists)",
			fields: fields{
				fmtProxy:     fmtproxy.New(),
				osProxy:      osproxy.New(),
				strconvProxy: strconvproxy.New(),
			},
			args: args{
				dirPath: testDirPath,
			},
			wantErr: false,
			setup: func() {
				if err := osProxy.MkdirAll(testDirPath, osProxy.FileMode(0755)); err != nil {
					t.Errorf("OsProxy.MkdirAll(%v) : error =\n%v", testDirPath, err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(testDirPath); err != nil {
					t.Errorf("OsProxy.RemoveAll(%v) : error =\n%v", testDirPath, err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(
				tt.fields.fmtProxy,
				tt.fields.osProxy,
				tt.fields.strconvProxy,
			)
			if tt.setup != nil {
				tt.setup()
			}
			if err := u.CreateDirIfNotExist(tt.args.dirPath); (err != nil) != tt.wantErr {
				t.Errorf("Utility.CreateDirIfNotExist() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
