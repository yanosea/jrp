package completion

import (
	"bytes"
	"testing"

	c "github.com/spf13/cobra"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

func TestNewCompletionZshCommand(t *testing.T) {
	type args struct {
		cobra  proxy.Cobra
		output *string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{
				cobra:  proxy.NewCobra(),
				output: new(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCompletionZshCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewCompletionZshCommand() = %v, want not nil", got)
			} else {
				if err := got.RunE(&c.Command{}, []string{}); err != nil {
					t.Errorf("Failed to run the completion zsh command: %v", err)
				}
			}
		})
	}
}

func Test_runCompletionZsh(t *testing.T) {
	var output string
	rootCmd := &c.Command{Use: "jrp"}
	subCmd := &c.Command{Use: "sub"}
	rootCmd.AddCommand(subCmd)
	buf := new(bytes.Buffer)
	if err := rootCmd.GenZshCompletion(buf); err != nil {
		t.Errorf("Failed to generate the bash completion: %v", err)
	}
	wantOutput := buf.String()

	type args struct {
		cmd    *c.Command
		output *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing",
			args: args{
				cmd:    rootCmd,
				output: &output,
			},
			want:    wantOutput,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCompletionZsh(tt.args.cmd, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runCompletionZsh() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != *tt.args.output {
				t.Errorf("runCompletionZsh() = %v, want %v", *tt.args.output, tt.want)
			}
		})
	}
}
