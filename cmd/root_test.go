package cmd

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	type args struct {
		globalOption *GlobalOption
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
		setup   func()
	}{
		{
			name:    "positive testing",
			args:    args{globalOption: &GlobalOption{Out: os.Stdout, ErrOut: os.Stderr}},
			want:    0,
			wantErr: false,
			setup:   nil,
		}, {
			// TODO : negative testing
			name:    "negative testing (rootCmd.Execute() fails)",
			args:    args{},
			want:    1,
			wantErr: true,
			setup: func() {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if got := tt.args.globalOption.Execute(); (got != 0) != tt.wantErr {
				t.Errorf("Execute() : got = %v, want = %v", got, tt.want)
			}
		})
	}
}
