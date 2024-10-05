package promptuiproxy

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New()
		})
	}
}

func TestPromptuiProxy_NewPrompt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.NewPrompt()
		})
	}
}
