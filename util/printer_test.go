package util

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPrintlnWithWriter(t *testing.T) {
	type args struct {
		message string
		writer  io.Writer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing (stdout)",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "test message stdout\n",
		}, {
			name: "positive testing (stderr)",
			args: args{
				message: "test message stderr",
				writer:  os.Stderr,
			},
			want: "test message stderr\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.args.writer = &buf
			PrintlnWithWriter(tt.args.writer, tt.args.message)
			if got := buf.String(); got != tt.want {
				t.Errorf("PrintlnWithWriter() : got =  %v, want = %v", got, tt.want)
			}
		})
	}
}
