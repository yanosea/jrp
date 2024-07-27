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
			origExitFunc := exitFunc
			exitFunc = func(code int) {
				if code != 0 {
					t.Fatalf("exitFunc was called with code %v", code)
				}
			}
			defer func() {
				exitFunc = origExitFunc
			}()
			main()
		})
	}
}
