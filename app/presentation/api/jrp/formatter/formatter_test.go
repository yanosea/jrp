package formatter

import (
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
			name: "positive testing (format is json)",
			args: args{
				format: "json",
			},
			want:    &JsonFormatter{},
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
