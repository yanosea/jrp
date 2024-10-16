package constant

import (
	"testing"
)

func TestGetHistorySearchAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"se", "S"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHistorySearchAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetHistorySearchAliases() : [%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
