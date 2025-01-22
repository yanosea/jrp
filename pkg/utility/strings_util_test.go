package utility

import (
	"reflect"
	"testing"
)

func TestNewStringsUtil(t *testing.T) {
	tests := []struct {
		name string
		want StringsUtil
	}{
		{
			name: "positive testing",
			want: &stringsUtil{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStringsUtil(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStringsUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringsUtil_RemoveNewLines(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		s    *stringsUtil
		args args
		want string
	}{
		{
			name: "positive testing (string with newlines)",
			s:    &stringsUtil{},
			args: args{
				str: "hello\nworld\n",
			},
			want: "helloworld",
		},
		{
			name: "positive testing (string without newlines)",
			s:    &stringsUtil{},
			args: args{
				str: "helloworld",
			},
			want: "helloworld",
		},
		{
			name: "positive testing (empty string)",
			s:    &stringsUtil{},
			args: args{
				str: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &stringsUtil{}
			if got := s.RemoveNewLines(tt.args.str); got != tt.want {
				t.Errorf("stringsUtil.RemoveNewLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringsUtil_RemoveSpaces(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		s    *stringsUtil
		args args
		want string
	}{
		{
			name: "positive testing (string with spaces)",
			s:    &stringsUtil{},
			args: args{
				str: "hello world",
			},
			want: "helloworld",
		},
		{
			name: "positive testing (string without spaces)",
			s:    &stringsUtil{},
			args: args{
				str: "helloworld",
			},
			want: "helloworld",
		},
		{
			name: "positive testing (empty string)",
			s:    &stringsUtil{},
			args: args{
				str: "",
			},
			want: "",
		},
		{
			name: "positive testing (multiple spaces)",
			s:    &stringsUtil{},
			args: args{
				str: "hello  world  ",
			},
			want: "helloworld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &stringsUtil{}
			if got := s.RemoveSpaces(tt.args.str); got != tt.want {
				t.Errorf("stringsUtil.RemoveSpaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringsUtil_RemoveTabs(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		s    *stringsUtil
		args args
		want string
	}{
		{
			name: "positive testing (string with tabs)",
			s:    &stringsUtil{},
			args: args{
				str: "hello\tworld",
			},
			want: "helloworld",
		},
		{
			name: "positive testing (string without tabs)",
			s:    &stringsUtil{},
			args: args{
				str: "helloworld",
			},
			want: "helloworld",
		},
		{
			name: "positive testing (empty string)",
			s:    &stringsUtil{},
			args: args{
				str: "",
			},
			want: "",
		},
		{
			name: "positive testing (multiple tabs)",
			s:    &stringsUtil{},
			args: args{
				str: "hello\t\tworld\t",
			},
			want: "helloworld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &stringsUtil{}
			if got := s.RemoveTabs(tt.args.str); got != tt.want {
				t.Errorf("stringsUtil.RemoveTabs() = %v, want %v", got, tt.want)
			}
		})
	}
}
