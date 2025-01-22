package formatter

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"
)

// TableFormatter is a struct that formats the output of jrp cli.
type TableFormatter struct{}

// NewTableFormatter returns a new instance of the TableFormatter struct.
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

var (
	// t is a variable to store the table writer with the default values for injecting the dependencies in testing.
	t = utility.NewTableWriterUtil(proxy.NewTableWriter())
)

// tableData is a struct that holds the data of a table.
type tableData struct {
	header []string
	rows   [][]string
}

// Format formats the output of jrp cli.
func (f *TableFormatter) Format(result interface{}) string {
	var data tableData

	switch v := result.(type) {
	case []*jrpApp.GenerateJrpUseCaseOutputDto:
		data = f.formatGenerateJrp(v)
	case []*jrpApp.GetHistoryUseCaseOutputDto:
		data = f.formatHistory(v, func(h interface{}) (int, string, string, string, int, time.Time, time.Time) {
			dto := h.(*jrpApp.GetHistoryUseCaseOutputDto)
			return dto.ID, dto.Phrase, dto.Prefix, dto.Suffix, dto.IsFavorited, dto.CreatedAt, dto.UpdatedAt
		})
	case []*jrpApp.SearchHistoryUseCaseOutputDto:
		data = f.formatHistory(v, func(h interface{}) (int, string, string, string, int, time.Time, time.Time) {
			dto := h.(*jrpApp.SearchHistoryUseCaseOutputDto)
			return dto.ID, dto.Phrase, dto.Prefix, dto.Suffix, dto.IsFavorited, dto.CreatedAt, dto.UpdatedAt
		})
	default:
		return ""
	}

	return f.getTableString(data)
}

// formatGenerateJrp formats the output of the GenerateJrp use case.
func (f *TableFormatter) formatGenerateJrp(items []*jrpApp.GenerateJrpUseCaseOutputDto) tableData {
	header := []string{"phrase", "prefix", "suffix", "created_at"}

	noId := slices.ContainsFunc(items, func(dto *jrpApp.GenerateJrpUseCaseOutputDto) bool {
		return dto.ID == 0
	})
	if !noId {
		header = append([]string{"id"}, header...)
	}

	var rows [][]string
	for _, jrp := range items {
		row := []string{jrp.Phrase, jrp.Prefix, jrp.Suffix, jrp.CreatedAt.Format("2006-01-02 15:04:05")}
		if !noId {
			row = append([]string{strconv.Itoa(jrp.ID)}, row...)
		}
		rows = append(rows, row)
	}

	if !noId {
		rows = f.addTotalRow(rows)
	}

	return tableData{header: header, rows: rows}
}

// formatHistory formats the output of the GetHistory and SearchHistory use cases.
func (f *TableFormatter) formatHistory(items interface{}, getData func(interface{}) (int, string, string, string, int, time.Time, time.Time)) tableData {
	header := []string{"id", "phrase", "prefix", "suffix", "is_favorited", "created_at", "updated_at"}
	var rows [][]string

	addRows := func(v interface{}) {
		switch v := v.(type) {
		case []*jrpApp.GetHistoryUseCaseOutputDto:
			for _, item := range v {
				id, phrase, prefix, suffix, isFavorited, createdAt, updatedAt := getData(item)
				favorited := ""
				if isFavorited == 1 {
					favorited = "○"
				}
				rows = append(rows, []string{
					strconv.Itoa(id),
					phrase,
					prefix,
					suffix,
					favorited,
					createdAt.Format("2006-01-02 15:04:05"),
					updatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		case []*jrpApp.SearchHistoryUseCaseOutputDto:
			for _, item := range v {
				id, phrase, prefix, suffix, isFavorited, createdAt, updatedAt := getData(item)
				favorited := ""
				if isFavorited == 1 {
					favorited = "○"
				}
				rows = append(rows, []string{
					strconv.Itoa(id),
					phrase,
					prefix,
					suffix,
					favorited,
					createdAt.Format("2006-01-02 15:04:05"),
					updatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		default:
			return
		}
	}
	addRows(items)

	if len(rows) <= 0 {
		return tableData{}
	}

	rows = f.addTotalRow(rows)
	return tableData{header: header, rows: rows}
}

// addTotalRow adds a total row to the table.
func (f *TableFormatter) addTotalRow(rows [][]string) [][]string {
	if len(rows) == 0 {
		return [][]string{}
	}

	emptyRow := make([]string, len(rows[0]))
	for i := range emptyRow {
		emptyRow[i] = ""
	}

	totalRow := make([]string, len(rows[0]))
	totalRow[0] = fmt.Sprintf("TOTAL : %d jrps!", len(rows))

	rows = append(rows, emptyRow)
	rows = append(rows, totalRow)

	return rows
}

// getTableString returns a string representation of a table.
func (f *TableFormatter) getTableString(data tableData) string {
	if len(data.header) == 0 || len(data.rows) == 0 {
		return ""
	}

	tableString := &strings.Builder{}
	table := t.GetNewDefaultTable(tableString)
	table.SetHeader(data.header)
	table.AppendBulk(data.rows)
	table.Render()
	return strings.TrimSuffix(tableString.String(), "\n")
}
