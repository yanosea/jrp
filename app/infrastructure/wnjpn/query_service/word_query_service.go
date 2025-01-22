package query_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/yanosea/jrp/app/application/wnjpn"
	"github.com/yanosea/jrp/app/infrastructure/database"
)

// WordQueryService is a struct that implements the WordQueryService interface.
type wordQueryService struct {
	connManager database.ConnectionManager
}

// NewWordQueryService returns a new instance of the WordQueryService struct.
func NewWordQueryService() wnjpn.WordQueryService {
	return &wordQueryService{
		connManager: database.GetConnectionManager(),
	}
}

// FindByLangAndPosIn is a method that fetches words by lang and pos.
func (w *wordQueryService) FindByLangIsAndPosIn(
	ctx context.Context,
	lang string,
	pos []string,
) ([]*wnjpn.FetchWordsDto, error) {
	var deferErr error
	conn, err := w.connManager.GetConnection(database.WNJpnDB)
	if err != nil {
		return nil, err
	}

	db, err := conn.Open()
	if err != nil {
		return nil, err
	}

	placeholders := make([]string, len(pos))
	for i := range pos {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(FindByLangIsAndPosInQuery, strings.Join(placeholders, ","))
	params := make([]interface{}, 0, len(pos)+1)
	params = append(params, lang)
	for _, p := range pos {
		params = append(params, p)
	}

	rows, err := db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	words := make([]*wnjpn.FetchWordsDto, 0)
	for rows.Next() {
		word := &wnjpn.FetchWordsDto{}
		if err := rows.Scan(
			&word.WordID,
			&word.Lang,
			&word.Lemma,
			&word.Pron,
			&word.Pos,
		); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, deferErr
}
