package jrp

import (
	"testing"

	"github.com/fatih/color"

	"github.com/yanosea/jrp/pkg/proxy"
)

func TestNewVersionCommand(t *testing.T) {
	type args struct {
		cobra   proxy.Cobra
		version string
		output  *string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{
				cobra:   proxy.NewCobra(),
				version: "0.0.0",
				output:  new(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVersionCommand(tt.args.cobra, tt.args.version, tt.args.output)
			if got == nil {
				t.Errorf("NewVersionCommand() = %v, want not nil", got)
			} else {
				if err := got.RunE(nil, []string{}); err != nil {
					t.Errorf("Failed to run the version command: %v", err)
				}
			}
		})
	}
}

func Test_runVersion(t *testing.T) {
	var output string
	origFormat := format

	type args struct {
		version string
		output  *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				version: "0.0.0",
				output:  &output,
			},
			want:    "jrp version 0.0.0",
			wantErr: false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (formatter.NewFormatter(\"plain\") failed)",
			args: args{
				version: "0.0.0",
				output:  &output,
			},
			want:    color.RedString("❌ Failed to create a formatter..."),
			wantErr: true,
			setup: func() {
				format = "test"
				output = ""
			},
			cleanup: func() {
				format = origFormat
				output = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := runVersion(tt.args.version, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *tt.args.output != tt.want {
				t.Errorf("runVersion() = %v, want %v", *tt.args.output, tt.want)
			}
		})
	}
}
