package utility

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

func TestNewTableWriterUtil(t *testing.T) {
	tableWriter := proxy.NewTableWriter()

	type args struct {
		tableWriter proxy.TableWriter
	}
	tests := []struct {
		name string
		args args
		want TableWriterUtil
	}{
		{
			name: "positive testing",
			args: args{
				tableWriter: tableWriter,
			},
			want: &tableWriterUtil{
				tableWriter: tableWriter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTableWriterUtil(tt.args.tableWriter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTableWriterUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tableWriterUtil_GetNewDefaultTable(t *testing.T) {
	tableWriter := proxy.NewTableWriter()

	type fields struct {
		TableWriter proxy.TableWriter
	}
	type args struct {
		writer io.Writer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive testing",
			fields: fields{
				TableWriter: tableWriter,
			},
			args: args{
				writer: os.Stdout,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &tableWriterUtil{
				tableWriter: tt.fields.TableWriter,
			}
			if got := tr.GetNewDefaultTable(tt.args.writer); got == nil {
				t.Errorf("tableWriterUtil.GetNewDefaultTable() = %v, want not nil", got)
			}
		})
	}
}
