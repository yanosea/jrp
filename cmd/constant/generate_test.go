package constant

import (
	"testing"
)

func TestGetGenerateAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"gen", "g"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGenerateAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetDownloadAliases()[%v] =\n%v, want\n%v", i, got, tt.want)
				}
			}
		})
	}
}
