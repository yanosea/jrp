package formatter

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

// TableFormatter is a struct that formats the output of jrp cli.
type TableFormatter struct{}

// NewTableFormatter returns a new instance of the TableFormatter struct.
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

var (
	// Tu is a variable to store the table writer with the default values for injecting the dependencies in testing.
	Tu = utility.NewTableWriterUtil(proxy.NewTableWriter())
)

// tableData is a struct that holds the data of a table.
type tableData struct {
	header []string
	rows   [][]string
}

// Format formats the output of jrp cli.
func (f *TableFormatter) Format(result any) (string, error) {
	var data tableData

	switch v := result.(type) {
	case []*jrpApp.GenerateJrpUseCaseOutputDto:
		data = f.formatGenerateJrp(v)
	case []*jrpApp.GetHistoryUseCaseOutputDto:
		data = f.formatHistory(v, func(h any) (int, string, string, string, int, time.Time, time.Time) {
			dto := h.(*jrpApp.GetHistoryUseCaseOutputDto)
			return dto.ID, dto.Phrase, dto.Prefix, dto.Suffix, dto.IsFavorited, dto.CreatedAt, dto.UpdatedAt
		})
	case []*jrpApp.SearchHistoryUseCaseOutputDto:
		data = f.formatHistory(v, func(h any) (int, string, string, string, int, time.Time, time.Time) {
			dto := h.(*jrpApp.SearchHistoryUseCaseOutputDto)
			return dto.ID, dto.Phrase, dto.Prefix, dto.Suffix, dto.IsFavorited, dto.CreatedAt, dto.UpdatedAt
		})
	default:
		return "", errors.New("invalid result type")
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

	return tableData{header: header, rows: rows}
}

// formatHistory formats the output of the GetHistory and SearchHistory use cases.
func (f *TableFormatter) formatHistory(items any, getData func(any) (int, string, string, string, int, time.Time, time.Time)) tableData {
	header := []string{"id", "phrase", "prefix", "suffix", "is_favorited", "created_at", "updated_at"}
	var rows [][]string

	addRows := func(v any) {
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

	return tableData{header: header, rows: rows}
}

// getTotalFooter returns a total footer of the table.
func (f *TableFormatter) getTotalFooter(rows [][]string) []string {
	if len(rows) == 0 {
		return nil
	}

	totalRow := make([]string, len(rows[0]))
	for i := range totalRow {
		totalRow[i] = ""
	}
	lastColIndex := len(totalRow) - 1
	totalRow[lastColIndex] = fmt.Sprintf("TOTAL : %d jrps!", len(rows))

	return totalRow
}

// getTableString returns a string representation of a table.
func (f *TableFormatter) getTableString(data tableData) (string, error) {
	if len(data.header) == 0 || len(data.rows) == 0 {
		return "", nil
	}

	tableString := &strings.Builder{}
	table := Tu.GetNewDefaultTable(tableString)
	table.Header(data.header)
	shouldHaveFooter := true
	if len(data.header) > 0 && data.header[0] != "id" {
		shouldHaveFooter = false
	}
	if shouldHaveFooter {
		footer := f.getTotalFooter(data.rows)
		if footer != nil {
			table.Footer(footer)
		}
	}

	if err := table.Bulk(data.rows); err != nil {
		return "", err
	}
	if err := table.Render(); err != nil {
		return "", err
	}

	return strings.TrimSuffix(tableString.String(), "\n"), nil
}
