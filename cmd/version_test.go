package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
)

func TestNewVersionCommand(t *testing.T) {
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
				Use:   constant.VERSION_USE,
				Short: constant.VERSION_SHORT,
				Long:  constant.VERSION_LONG,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newVersionCommand(tt.args.globalOption)
			if got.Use != tt.want.Use || got.Short != tt.want.Short || got.Long != tt.want.Long {
				t.Errorf("newVersionCommand() : got = %v, want = %v", got, tt.want)
			}
			if err := got.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("newVersionCommand().Execute() : error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestGlobalOption_version(t *testing.T) {
	type fields struct {
		Out    io.Writer
		ErrOut io.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "positive testing",
			fields: fields{
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GlobalOption{
				Out:    tt.fields.Out,
				ErrOut: tt.fields.ErrOut,
			}
			if err := g.version(); (err != nil) != tt.wantErr {
				t.Errorf("GlobalOption.version() : error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
