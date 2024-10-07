package constant

import (
	"testing"
)

func TestGetFavoriteAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"fav", "f"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFavoriteAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetFavoriteAliases()[%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
