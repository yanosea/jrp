package utility

import (
	"io"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// TableWriterUtil is an interface that contains the utility functions for writing tables.
type TableWriterUtil interface {
	GetNewDefaultTable(writer io.Writer) proxy.Table
}

// tableWriterUtil is a struct that contains the utility functions for writing tables.
type tableWriterUtil struct {
	tableWriter proxy.TableWriter
}

// NewTableWriterUtil returns a new instance of the TableWriterUtil interface.
func NewTableWriterUtil(tableWriter proxy.TableWriter) TableWriterUtil {
	return &tableWriterUtil{
		tableWriter: tableWriter,
	}
}

// GetNewDefaultTable returns a new instance of the default table.
func (t *tableWriterUtil) GetNewDefaultTable(writer io.Writer) proxy.Table {
	table := t.tableWriter.NewTable(writer)
	return table
}
