package util

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestFormatIndent(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message",
			},
			want: "  test message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatIndent(tt.args.message); got != tt.want {
				t.Errorf("FormatIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name: "positive testing",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "test message stdout\n",
		}, {
			name: "positive testing",
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
				t.Errorf("PrintlnWithWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintWithWriterWithBlankLineBelow(t *testing.T) {
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
			name: "positive testing",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "test message stdout\n\n",
		}, {
			name: "positive testing",
			args: args{
				message: "test message stderr",
				writer:  os.Stderr,
			},
			want: "test message stderr\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.args.writer = &buf
			PrintWithWriterWithBlankLineBelow(tt.args.writer, tt.args.message)
			if got := buf.String(); got != tt.want {
				t.Errorf("PrintWithWriterWithBlankLineBelow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintWithWriterWithBlankLineAbove(t *testing.T) {
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
			name: "positive testing",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "\ntest message stdout\n",
		}, {
			name: "positive testing",
			args: args{
				message: "test message stderr",
				writer:  os.Stderr,
			},
			want: "\ntest message stderr\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.args.writer = &buf
			PrintWithWriterWithBlankLineAbove(tt.args.writer, tt.args.message)
			if got := buf.String(); got != tt.want {
				t.Errorf("PrintWithWriterWithBlankLineAbove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintWithWriterBetweenBlankLine(t *testing.T) {
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
			name: "positive testing",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "\ntest message stdout\n\n",
		}, {
			name: "positive testing",
			args: args{
				message: "test message stderr",
				writer:  os.Stderr,
			},
			want: "\ntest message stderr\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.args.writer = &buf
			PrintWithWriterBetweenBlankLine(tt.args.writer, tt.args.message)
			if got := buf.String(); got != tt.want {
				t.Errorf("PrintWithWriterWithBlankLineAbove() = %v, want %v", got, tt.want)
			}
		})
	}
}
