package constant

import (
	"testing"
)

func TestGetInteractiveAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"int", "i"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInteractiveAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetInteractiveAliases() : [%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
