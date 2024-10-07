package constant

import (
	"testing"
)

func TestGetHistoryAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"hist", "h"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHistoryAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetHistoryAliases()[%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
