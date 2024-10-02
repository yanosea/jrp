package constant

import (
	"testing"
)

func TestGetFavoriteAddAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"ad", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFavoriteAddAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetFavoriteAddAliases()[%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
