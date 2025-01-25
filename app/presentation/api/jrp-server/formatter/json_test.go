package formatter

import (
	"errors"
	"reflect"
	"testing"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewJsonFormatter(t *testing.T) {
	tests := []struct {
		name string
		want *JsonFormatter
	}{
		{
			name: "positive testing",
			want: &JsonFormatter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsonFormatter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonFormatter_Format(t *testing.T) {
	origJu := Ju

	type args struct {
		result interface{}
	}
	tests := []struct {
		name    string
		f       *JsonFormatter
		args    args
		want    []byte
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing (result is *jrpApp.GenerateJrpUseCaseOutputDto)",
			args: args{
				result: &jrpApp.GenerateJrpUseCaseOutputDto{
					Phrase: "test",
				},
			},
			want:    []byte(`{"phrase":"test"}`),
			wantErr: false,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "negative testing (result is *jrpApp.GenerateJrpUseCaseOutputDto, Ju.Marshal(jjoDto) failed)",
			args: args{
				result: &jrpApp.GenerateJrpUseCaseOutputDto{
					Phrase: "test",
				},
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockJson := proxy.NewMockJson(mockCtrl)
				mockJson.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("JsonUtil.Marshal() failed"))
				Ju = utility.NewJsonUtil(mockJson)
			},
			cleanup: func() {
				Ju = origJu
			},
		},
		{
			name: "negative testing (result is invalid)",
			args: args{
				result: "invalid",
			},
			want:    nil,
			wantErr: true,
			setup:   nil,
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			f := &JsonFormatter{}
			got, err := f.Format(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
