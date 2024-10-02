package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strconv"
)

func TestNewCompletionBashCommand(t *testing.T) {
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
			got := NewCompletionBashCommand(tt.args.g)
			if err := got.Execute(); (err != nil) != tt.wantError {
				t.Errorf("NewCompletionBashCommand().Execute() : error =\n%v", err)
			}
		})
	}
}

func TestGlobalOption_completionBashRunE(t *testing.T) {
	globalOption := NewGlobalOption(fmtproxy.New(), osproxy.New(), strconvproxy.New())
	cmd := NewRootCommand(globalOption.Out, globalOption.ErrOut, globalOption.Args)

	type fields struct {
		Out            ioproxy.WriterInstanceInterface
		ErrOut         ioproxy.WriterInstanceInterface
		Args           []string
		Utility        utility.UtilityInterface
		NewRootCommand func(ow, ew ioproxy.WriterInstanceInterface, args []string) cobraproxy.CommandInstanceInterface
	}
	type args struct {
		c   *cobra.Command
		in0 []string
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
				c:   cmd.GetCommand(),
				in0: nil,
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
			if err := g.completionBashRunE(tt.args.c, tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("GlobalOption.completionBashRunE() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
		})
	}
}

func TestGlobalOption_completionBash(t *testing.T) {
	globalOption := NewGlobalOption(fmtproxy.New(), osproxy.New(), strconvproxy.New())
	cmd := NewRootCommand(globalOption.Out, globalOption.ErrOut, globalOption.Args)

	type fields struct {
		Out            ioproxy.WriterInstanceInterface
		ErrOut         ioproxy.WriterInstanceInterface
		Args           []string
		Utility        utility.UtilityInterface
		NewRootCommand func(ow, ew ioproxy.WriterInstanceInterface, args []string) cobraproxy.CommandInstanceInterface
	}
	type args struct {
		c *cobra.Command
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
				c: cmd.GetCommand(),
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
			if err := g.completionBash(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GlobalOption.completionBash() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
		})
	}
}
