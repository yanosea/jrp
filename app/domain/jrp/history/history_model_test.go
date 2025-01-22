package history

import (
	"database/sql"
	"reflect"
	"testing"
	"time"
)

func TestNewHistory(t *testing.T) {
	now := time.Now()
	type args struct {
		phrase      string
		prefix      string
		suffix      string
		isFavorited int
		createdAt   time.Time
		updatedAt   time.Time
	}
	tests := []struct {
		name string
		args args
		want *History
	}{
		{
			name: "positive testing",
			args: args{
				phrase:      "prefix test",
				prefix:      "prefix",
				suffix:      "",
				isFavorited: 1,
				createdAt:   now,
				updatedAt:   now,
			},
			want: &History{
				Phrase:      "prefix test",
				Prefix:      sql.NullString{String: "prefix", Valid: true},
				Suffix:      sql.NullString{String: "", Valid: false},
				IsFavorited: 1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHistory(tt.args.phrase, tt.args.prefix, tt.args.suffix, tt.args.isFavorited, tt.args.createdAt, tt.args.updatedAt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}
