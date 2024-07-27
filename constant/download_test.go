package constant

import (
	"testing"
)

func TestGetDownloadAliases(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "positive testing",
			want: []string{"dl", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDownloadAliases()
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetDownloadAliases()[%v] = %v, want %v", i, got, tt.want)
				}
			}
		})
	}
}
