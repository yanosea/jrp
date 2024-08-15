package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
)

func TestNewCompletionCommand(t *testing.T) {
	type args struct {
		globalOption *GlobalOption
	}
	tests := []struct {
		name    string
		args    args
		want    *cobra.Command
		wantErr bool
	}{
		{
			name: "positive testing",
			args: args{globalOption: &GlobalOption{Out: os.Stdout, ErrOut: os.Stderr}},
			want: &cobra.Command{
				Use:   constant.COMPLETION_USE,
				Short: constant.COMPLETION_SHORT,
				Long:  constant.COMPLETION_LONG,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCompletionCommand(tt.args.globalOption)
			if got.Use != tt.want.Use || got.Short != tt.want.Short || got.Long != tt.want.Long {
				t.Errorf("newCompletionCommand() : got = %v, want = %v", got, tt.want)
			}
			if err := got.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("newCompletionCommand().Execute() : error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
