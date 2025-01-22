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
	AppendBulk(rows [][]string)
	Render()
	SetAlignment(align int)
	SetAutoFormatHeaders(auto bool)
	SetAutoWrapText(auto bool)
	SetBorder(border bool)
	SetCenterSeparator(sep string)
	SetColumnSeparator(sep string)
	SetHeader(keys []string)
	SetHeaderAlignment(hAlign int)
	SetHeaderLine(line bool)
	SetNoWhiteSpace(allow bool)
	SetRowSeparator(sep string)
	SetTablePadding(padding string)
}

// tableProxy is a proxy struct that implements the Table interface.
type tableProxy struct {
	table *tablewriter.Table
}

// AppendBulk appends multiple rows to the table.
func (t *tableProxy) AppendBulk(rows [][]string) {
	t.table.AppendBulk(rows)
}

// Render renders the table.
func (t *tableProxy) Render() {
	t.table.Render()
}

// SetAlignment sets the alignment of the table.
func (t *tableProxy) SetAlignment(align int) {
	t.table.SetAlignment(align)
}

// SetAutoFormatHeaders sets the auto format headers of the table.
func (t *tableProxy) SetAutoFormatHeaders(auto bool) {
	t.table.SetAutoFormatHeaders(auto)
}

// SetAutoWrapText sets the auto wrap text of the table.
func (t *tableProxy) SetAutoWrapText(auto bool) {
	t.table.SetAutoWrapText(auto)
}

// SetBorder sets the border of the table.
func (t *tableProxy) SetBorder(border bool) {
	t.table.SetBorder(border)
}

// SetCenterSeparator sets the center separator of the table.
func (t *tableProxy) SetCenterSeparator(sep string) {
	t.table.SetCenterSeparator(sep)
}

// SetColumnSeparator sets the column separator of the table.
func (t *tableProxy) SetColumnSeparator(sep string) {
	t.table.SetColumnSeparator(sep)
}

// SetHeader sets the header of the table.
func (t *tableProxy) SetHeader(keys []string) {
	t.table.SetHeader(keys)
}

// SetHeaderAlignment sets the header alignment of the table.
func (t *tableProxy) SetHeaderAlignment(hAlign int) {
	t.table.SetHeaderAlignment(hAlign)
}

// SetHeaderLine sets the header line of the table.
func (t *tableProxy) SetHeaderLine(line bool) {
	t.table.SetHeaderLine(line)
}

// SetNoWhiteSpace sets the no white space of the table.
func (t *tableProxy) SetNoWhiteSpace(allow bool) {
	t.table.SetNoWhiteSpace(allow)
}

// SetRowSeparator sets the row separator of the table.
func (t *tableProxy) SetRowSeparator(sep string) {
	t.table.SetRowSeparator(sep)
}

// SetTablePadding sets the table padding of the table.
func (t *tableProxy) SetTablePadding(padding string) {
	t.table.SetTablePadding(padding)
}
