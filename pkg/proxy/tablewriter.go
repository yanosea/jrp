package proxy

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// TableWriter is an interface that provides a proxy of the methods of tablewriter.
type TableWriter interface {
	NewTable(writer io.Writer) Table
}

// tableWriterProxy is a proxy struct that implements the TableWriter interface.
type tableWriterProxy struct{}

// NewTableWriter returns a new instance of the TableWriter interface.
func NewTableWriter() TableWriter {
	return &tableWriterProxy{}
}

// NewTable returns a new instance of the tablewriter.Table.
func (t *tableWriterProxy) NewTable(writer io.Writer) Table {
	return &tableProxy{
		table: tablewriter.NewWriter(writer),
	}
}

// Table is an interface that provides a proxy of the methods of tablewriter.Table.
type Table interface {
	Bulk(rows [][]string) error
	Footer(elements []string)
	Render() error
	Header(elements []string)
}

// tableProxy is a proxy struct that implements the Table interface.
type tableProxy struct {
	table *tablewriter.Table
}

// Bulk appends multiple rows to the table.
func (t *tableProxy) Bulk(rows [][]string) error {
	return t.table.Bulk(rows)
}

// Footer sets the footer of the table.
func (t *tableProxy) Footer(elements []string) {
	t.table.Footer(elements)
}

// Render renders the table.
func (t *tableProxy) Render() error {
	return t.table.Render()
}

// Header sets the header of the table.
func (t *tableProxy) Header(elements []string) {
	t.table.Header(elements)
}
