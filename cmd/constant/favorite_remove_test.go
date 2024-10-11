package constant

import (
	"testing"
)

func TestGetFavoriteRemoveAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"rm", "r"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFavoriteRemoveAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetFavoriteRemoveAliases() : [%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
