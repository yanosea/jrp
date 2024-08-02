package main

import (
	"testing"
)

func TestMain(t *testing.T) {

	type want struct {
		exitCode int
	}
	tests := []struct {
		name  string
		want  want
		setup func()
	}{
		{
			name: "positive testing",
			want: want{exitCode: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			origOsExit := osExit
			osExit = func(code int) {
				if code != tt.want.exitCode {
					t.Errorf("main() : exit code = %v, want = %v", code, tt.want.exitCode)
				}
			}
			defer func() {
				osExit = origOsExit
			}()

			main()
		})
	}
}
