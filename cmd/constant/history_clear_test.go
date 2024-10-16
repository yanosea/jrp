package constant

import (
	"testing"
)

func TestGetHistoryClearAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"cl", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHistoryClearAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetHistoryClearAliases() : [%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
