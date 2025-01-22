package formatter

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		want    Formatter
		wantErr bool
	}{
		{
			name: "positive testing (format is plain)",
			args: args{
				format: "plain",
			},
			want:    &PlainFormatter{},
			wantErr: false,
		},
		{
			name: "positive testing (format is table)",
			args: args{
				format: "table",
			},
			want:    &TableFormatter{},
			wantErr: false,
		},
		{
			name: "negative testing (format is invalid)",
			args: args{
				format: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFormatter(tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendErrorToOutput(t *testing.T) {
	type args struct {
		err    error
		output string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing (err is nil, output is empty)",
			args: args{
				err:    nil,
				output: "",
			},
			want: "",
		},
		{
			name: "positive testing (err is not nil, output is empty)",
			args: args{
				err:    errors.New("test"),
				output: "",
			},
			want: Red("Error : test"),
		},
		{
			name: "positive testing (err is nil, output is not empty)",
			args: args{
				err:    nil,
				output: "test",
			},
			want: Red("test"),
		},
		{
			name: "positive testing (err is not nil, output is not empty)",
			args: args{
				err:    errors.New("test"),
				output: "test",
			},
			want: Red("test\nError : test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendErrorToOutput(tt.args.err, tt.args.output); got != tt.want {
				t.Errorf("AppendErrorToOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
