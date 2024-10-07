package jrpwriter

import (
	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/tablewriter"
)

// JrpWritable is an interface for JrpWriter.
type JrpWritable interface {
	WriteGenerateResultAsTable(writer ioproxy.WriterInstanceInterface, jrps []model.Jrp)
	WriteAsTable(writer ioproxy.WriterInstanceInterface, jrps []model.Jrp)
}

// JrpWriter is a struct that implements JrpWritable.
type JrpWriter struct {
	StrconvProxy     strconvproxy.Strconv
	TableWriterProxy tablewriterproxy.TableWriter
}

// New is a constructor for JrpWriter.
func New(
	strconvProxy strconvproxy.Strconv,
	tableWriterProxy tablewriterproxy.TableWriter,
) *JrpWriter {
	return &JrpWriter{
		StrconvProxy:     strconvProxy,
		TableWriterProxy: tableWriterProxy,
	}
}

// WriteGenerateResultAsTable writes the generate result as table.
func (j *JrpWriter) WriteGenerateResultAsTable(writer ioproxy.WriterInstanceInterface, jrps []model.Jrp) {
	if jrps == nil || len(jrps) <= 0 {
		return
	}

	headers := []string{"phrase", "prefix", "suffix", "created_at"}
	rowFunc := func(jrp model.Jrp) []string {
		prefix := ""
		if jrp.Prefix.FieldNullString.Valid {
			prefix = jrp.Prefix.FieldNullString.String
		}
		suffix := ""
		if jrp.Suffix.FieldNullString.Valid {
			suffix = jrp.Suffix.FieldNullString.String
		}
		return []string{
			jrp.Phrase,
			prefix,
			suffix,
			jrp.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	j.writeTable(writer, jrps, headers, rowFunc)
}

// WriteAsTable writes the jrps as table.
func (j *JrpWriter) WriteAsTable(writer ioproxy.WriterInstanceInterface, jrps []model.Jrp) {
	if jrps == nil || len(jrps) <= 0 {
		return
	}

	headers := []string{"id", "phrase", "prefix", "suffix", "is_favorited", "created_at", "updated_at"}
	rowFunc := func(jrp model.Jrp) []string {
		prefix := ""
		if jrp.Prefix.FieldNullString.Valid {
			prefix = jrp.Prefix.FieldNullString.String
		}
		suffix := ""
		if jrp.Suffix.FieldNullString.Valid {
			suffix = jrp.Suffix.FieldNullString.String
		}
		isFavorited := ""
		if jrp.IsFavorited == 1 {
			isFavorited = "○"
		}
		return []string{
			j.StrconvProxy.Itoa(jrp.ID),
			jrp.Phrase,
			prefix,
			suffix,
			isFavorited,
			jrp.CreatedAt.Format("2006-01-02 15:04:05"),
			jrp.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	j.writeTable(writer, jrps, headers, rowFunc)
}

// writeTable writes the table.
func (j *JrpWriter) writeTable(writer ioproxy.WriterInstanceInterface, jrps []model.Jrp, headers []string, rowFunc func(model.Jrp) []string) {
	if jrps == nil || len(jrps) <= 0 {
		return
	}

	var rows [][]string
	for _, jrp := range jrps {
		rows = append(rows, rowFunc(jrp))
	}
	total := j.StrconvProxy.Itoa(len(rows))
	rows = append(rows, make([]string, len(headers)))
	rows = append(rows, append([]string{"TOTAL : " + total}, make([]string, len(headers)-1)...))

	table := j.getDefaultTableWriter(writer)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}

// getDefaultTableWriter gets the default table instance.
func (j *JrpWriter) getDefaultTableWriter(o ioproxy.WriterInstanceInterface) tablewriterproxy.TableInstanceInterface {
	table := j.TableWriterProxy.NewTable(o)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriterproxy.ALIGN_LEFT)
	table.SetAlignment(tablewriterproxy.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	return table
}
