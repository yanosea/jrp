package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yanosea/jrp/constant"
)

func Test_newCompletionBashCommand(t *testing.T) {
	type args struct {
		globalOption *GlobalOption
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		{
			name: "positive testing",
			args: args{globalOption: &GlobalOption{Out: os.Stdout, ErrOut: os.Stderr}},
			want: &cobra.Command{
				Use:   constant.COMPLETION_BASH_USE,
				Short: constant.COMPLETION_BASH_SHORT,
				Long:  constant.COMPLETION_BASH_LONG,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCompletionBashCommand(tt.args.globalOption)
			if got.Use != tt.want.Use || got.Short != tt.want.Short || got.Long != tt.want.Long {
				t.Errorf("newCompletionBashCommand() : got = %v, want = %v", got, tt.want)
			}
		})
	}
}
