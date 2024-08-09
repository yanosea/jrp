package test

import (
	"fmt"
	"os"
	"testing"
)

func TestCaptureOutput(t *testing.T) {
	type args struct {
		t   *testing.T
		fnc func()
	}
	tests := []struct {
		name       string
		args       args
		wantStdOut string
		wantStdErr string
	}{
		{
			name: "positive testing",
			args: args{
				t: t,
				fnc: func() {
					fmt.Println("stdout")
					fmt.Fprintln(os.Stderr, "stderr")
				},
			},
			wantStdOut: "stdout\n",
			wantStdErr: "stderr\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdOut, gotStdErr := CaptureOutput(tt.args.t, tt.args.fnc)
			if gotStdOut != tt.wantStdOut {
				t.Errorf("CaptureOutput() : gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("CaptureOutput() : gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}
