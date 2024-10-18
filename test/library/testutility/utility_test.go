package testutility

import "testing"

func TestReplaceDoubleSlashToSingleSlash(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				input: "//path/to/file",
			},
			want: "/path/to/file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceDoubleSlashToSingleSlash(tt.args.input); got != tt.want {
				t.Errorf("ReplaceDoubleSlashToSingleSlash() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestRemoveTabAndSpaceAndLf(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				input: "test  \ttest\ntest",
			},
			want: "testtesttest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveTabAndSpaceAndLf(tt.args.input); got != tt.want {
				t.Errorf("RemoveTabAndSpaceAndLf() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}
