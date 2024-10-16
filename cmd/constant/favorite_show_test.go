package constant

import (
	"testing"
)

func TestGetFavoriteShowAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"sh", "s"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFavoriteShowAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetFavoriteShowAliases : ()[%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
