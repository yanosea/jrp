package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origOsExit := osExit
			osExit = func(code int) {
				if code != 0 {
					t.Fatalf("osExit was called with code %v", code)
				}
			}
			defer func() {
				osExit = origOsExit
			}()
			main()
		})
	}
}
