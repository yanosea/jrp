package formatter

import (
	"reflect"
	"testing"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
)

func TestNewPlainFormatter(t *testing.T) {
	tests := []struct {
		name string
		want *PlainFormatter
	}{
		{
			name: "positive testing",
			want: &PlainFormatter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlainFormatter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlainFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlainFormatter_Format(t *testing.T) {
	type args struct {
		result interface{}
	}
	tests := []struct {
		name string
		f    *PlainFormatter
		args args
		want string
	}{
		{
			name: "positive testing (result is *jrpApp.GetVersionUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: &jrpApp.GetVersionUseCaseOutputDto{
					Version: "0.0.0",
				},
			},
			want: "jrp version 0.0.0",
		},
		{
			name: "positive testing (result is []*jrpApp.GenerateJrpUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*jrpApp.GenerateJrpUseCaseOutputDto{
					{
						Phrase: "phrase1",
					},
					{
						Phrase: "phrase2",
					},
				},
			},
			want: "phrase1\nphrase2",
		},
		{
			name: "positive testing (result is []*jrpApp.GetHistoryUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*jrpApp.GetHistoryUseCaseOutputDto{
					{
						Phrase: "phrase1",
					},
					{
						Phrase: "phrase2",
					},
				},
			},
			want: "phrase1\nphrase2",
		},
		{
			name: "positive testing (result is []*jrpApp.SearchHistoryUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*jrpApp.SearchHistoryUseCaseOutputDto{
					{
						Phrase: "phrase1",
					},
					{
						Phrase: "phrase2",
					},
				},
			},
			want: "phrase1\nphrase2",
		},
		{
			name: "negative testing (result is invalid)",
			f:    &PlainFormatter{},
			args: args{
				result: "invalid",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PlainFormatter{}
			if got := f.Format(tt.args.result); got != tt.want {
				t.Errorf("PlainFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
