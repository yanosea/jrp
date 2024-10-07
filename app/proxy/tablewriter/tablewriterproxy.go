package tablewriterproxy

import (
	"github.com/olekukonko/tablewriter"

	"github.com/yanosea/jrp/app/proxy/io"
)

// TableWriter is an interface for tablewriter.
type TableWriter interface {
	NewTable(writer ioproxy.WriterInstanceInterface) *TableInstance
}

// TableWriterProxy is a struct that implements TableWriter.
type TableWriterProxy struct{}

// New is a constructor for TableWriterProxy.
func New() TableWriter {
	return &TableWriterProxy{}
}

// NewTable is a proxy for tablewriter.NewTable.
func (*TableWriterProxy) NewTable(writer ioproxy.WriterInstanceInterface) *TableInstance {
	return &TableInstance{FieldTable: tablewriter.NewWriter(writer)}
}
