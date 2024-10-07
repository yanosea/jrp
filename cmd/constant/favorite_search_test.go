package constant

import (
	"testing"
)

func TestGetFavoriteSearchAliases(t *testing.T) {
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
			got := GetFavoriteSearchAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetFavoriteSearchAliases()[%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
