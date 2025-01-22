package wnjpn

import (
	"context"
	"database/sql"
)

// FetchWordsDto is a DTO struct that contains the input data of the FetchWordsUseCase.
type FetchWordsDto struct {
	WordID int
	Lang   sql.NullString
	Lemma  sql.NullString
	Pron   sql.NullString
	Pos    sql.NullString
}

// WordQueryService is an interface that provides the methods to query the words.
type WordQueryService interface {
	FindByLangIsAndPosIn(ctx context.Context, lang string, pos []string) ([]*FetchWordsDto, error)
}
