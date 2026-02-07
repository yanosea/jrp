package proxy

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// TableWriter is an interface that provides a proxy of the methods of tablewriter.
type TableWriter interface {
	NewTable(writer io.Writer, opts ...tablewriter.Option) Table
}

// tableWriterProxy is a proxy struct that implements the TableWriter interface.
type tableWriterProxy struct{}

// NewTableWriter returns a new instance of the TableWriter interface.
func NewTableWriter() TableWriter {
	return &tableWriterProxy{}
}

// NewTable returns a new instance of the tablewriter.Table.
func (t *tableWriterProxy) NewTable(writer io.Writer, opts ...tablewriter.Option) Table {
	return &tableProxy{
		table: tablewriter.NewTable(writer, opts...),
	}
}

// Table is an interface that provides a proxy of the methods of tablewriter.Table.
type Table interface {
	Header(elements ...any)
	Bulk(rows [][]string) error
	Render() error
}

// tableProxy is a proxy struct that implements the Table interface.
type tableProxy struct {
	table *tablewriter.Table
}

// Header sets the header of the table.
func (t *tableProxy) Header(elements ...any) {
	t.table.Header(elements...)
}

// Bulk appends multiple rows to the table.
func (t *tableProxy) Bulk(rows [][]string) error {
	return t.table.Bulk(rows)
}

// Render renders the table.
func (t *tableProxy) Render() error {
	return t.table.Render()
}
