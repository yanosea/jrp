package tablewriterproxy

import (
	"github.com/olekukonko/tablewriter"
)

// TableInstanceInterface is an interface for tablewriter.Table.
type TableInstanceInterface interface {
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

// TableInstance is a struct that implements TableInstanceInterface.
type TableInstance struct {
	FieldTable *tablewriter.Table
}

// AppendBulk is a proxy for tablewriter.Table.AppendBulk.
func (t *TableInstance) AppendBulk(rows [][]string) {
	t.FieldTable.AppendBulk(rows)
}

// Render is a proxy for tablewriter.Table.Render.
func (t *TableInstance) Render() {
	t.FieldTable.Render()
}

// SetAlignment is a proxy for tablewriter.Table.SetAlignment.
func (t *TableInstance) SetAlignment(align int) {
	t.FieldTable.SetAlignment(align)
}

// SetAutoFormatHeaders is a proxy for tablewriter.Table.SetAutoFormatHeaders.
func (t *TableInstance) SetAutoFormatHeaders(auto bool) {
	t.FieldTable.SetAutoFormatHeaders(auto)
}

// SetAutoWrapText is a proxy for tablewriter.Table.SetAutoWrapText.
func (t *TableInstance) SetAutoWrapText(auto bool) {
	t.FieldTable.SetAutoWrapText(auto)
}

// SetBorder is a proxy for tablewriter.Table.SetBorder.
func (t *TableInstance) SetBorder(border bool) {
	t.FieldTable.SetBorder(border)
}

// SetCenterSeparator is a proxy for tablewriter.Table.SetCenterSeparator.
func (t *TableInstance) SetCenterSeparator(sep string) {
	t.FieldTable.SetCenterSeparator(sep)
}

// SetColumnSeparator is a proxy for tablewriter.Table.SetColumnSeparator.
func (t *TableInstance) SetColumnSeparator(sep string) {
	t.FieldTable.SetColumnSeparator(sep)
}

// SetHeader is a proxy for tablewriter.Table.SetHeader.
func (t *TableInstance) SetHeader(keys []string) {
	t.FieldTable.SetHeader(keys)
}

// SetHeaderAlignment is a proxy for tablewriter.Table.SetHeaderAlignment.
func (t *TableInstance) SetHeaderAlignment(hAlign int) {
	t.FieldTable.SetHeaderAlignment(hAlign)
}

// SetHeaderLine is a proxy for tablewriter.Table.SetHeaderLine.
func (t TableInstance) SetHeaderLine(line bool) {
	t.FieldTable.SetHeaderLine(line)
}

// SetNoWhiteSpace is a proxy for tablewriter.Table.SetNoWhiteSpace.
func (t *TableInstance) SetNoWhiteSpace(allow bool) {
	t.FieldTable.SetNoWhiteSpace(allow)
}

// SetRowSeparator is a proxy for tablewriter.Table.SetRowSeparator.
func (t *TableInstance) SetRowSeparator(sep string) {
	t.FieldTable.SetRowSeparator(sep)
}

// SetTablePadding is a proxy for tablewriter.Table.SetTablePadding.
func (t *TableInstance) SetTablePadding(padding string) {
	t.FieldTable.SetTablePadding(padding)
}
