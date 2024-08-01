package cmd

import (
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "positive testing",
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Execute(); (got != 0) != tt.wantErr {
				t.Errorf("Execute() : got = %v, want = %v", got, tt.want)
			}
		})
	}
}
