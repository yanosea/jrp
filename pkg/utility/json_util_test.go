package utility

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

func TestNewJsonUtil(t *testing.T) {
	json := proxy.NewJson()

	type args struct {
		json proxy.Json
	}
	tests := []struct {
		name string
		args args
		want JsonUtil
	}{
		{
			name: "positive testing",
			args: args{
				json: json,
			},
			want: &jsonUtil{
				json: json,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsonUtil(tt.args.json); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUtil_Marshal(t *testing.T) {
	type fields struct {
		json proxy.Json
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "positive testing",
			fields: fields{
				json: proxy.NewJson(),
			},
			args: args{
				v: map[string]interface{}{},
			},
			want:    []byte("{}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ju := &jsonUtil{
				json: tt.fields.json,
			}
			got, err := ju.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonUtil.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonUtil.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
