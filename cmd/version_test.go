package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strconv"

	"github.com/yanosea/jrp/test/testutility"
)

func TestNewVersionCommand(t *testing.T) {
	type args struct {
		g *GlobalOption
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name: "positive testing",
			args: args{
				g: NewGlobalOption(
					fmtproxy.New(),
					osproxy.New(),
					strconvproxy.New(),
				),
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVersionCommand(tt.args.g)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewVersionCommand().Execute() : error =\n%v", err)
			}
		})
	}
}

func TestGlobalOption_versionRunE(t *testing.T) {
	globalOption := NewGlobalOption(fmtproxy.New(), osproxy.New(), strconvproxy.New())

	type fields struct {
		Out            ioproxy.WriterInstanceInterface
		ErrOut         ioproxy.WriterInstanceInterface
		Args           []string
		Utility        utility.UtilityInterface
		NewRootCommand func(ow, ew ioproxy.WriterInstanceInterface, args []string) cobraproxy.CommandInstanceInterface
	}
	type args struct {
		in0 *cobra.Command
		in1 []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:            globalOption.Out,
				ErrOut:         globalOption.ErrOut,
				Args:           globalOption.Args,
				Utility:        globalOption.Utility,
				NewRootCommand: globalOption.NewRootCommand,
			},
			args: args{
				in0: nil,
				in1: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GlobalOption{
				Out:            tt.fields.Out,
				ErrOut:         tt.fields.ErrOut,
				Args:           tt.fields.Args,
				Utility:        tt.fields.Utility,
				NewRootCommand: tt.fields.NewRootCommand,
			}
			if err := g.versionRunE(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("GlobalOption.versionRunE() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
		})
	}
}

func TestGlobalOption_version(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	globalOption := NewGlobalOption(fmtproxy.New(), osProxy, strconvproxy.New())
	globalOption.Out = capturer.OutBuffer
	globalOption.ErrOut = capturer.ErrBuffer

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
			name: "positive testing",
			fields: fields{
				t: t,
				fnc: func() {
					globalOption.version()
				},
				capturer: capturer,
			},
			wantStdOut: "jrp version devel\n",
			wantStdErr: "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			if err != nil {
				t.Errorf("GlobalOption.completion() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if stdout != tt.wantStdOut {
				t.Errorf("GlobalOption.completion() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("GlobalOption.completion() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}
