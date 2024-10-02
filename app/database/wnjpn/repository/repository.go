package repository

import (
	"github.com/yanosea/jrp/app/database/wnjpn/model"
	"github.com/yanosea/jrp/app/database/wnjpn/repository/query"
	"github.com/yanosea/jrp/app/proxy/sql"
)

// WNJpnRepositoryInterface is an interface for WNJpnRepository.
type WNJpnRepositoryInterface interface {
	GetAllAVNWords(wnJpnDBFilePath string) ([]model.Word, error)
	GetAllNWords(wnJpnDBFilePath string) ([]model.Word, error)
	GetAllAVWords(wnJpnDBFilePath string) ([]model.Word, error)
}

// WNJpnRepository is a struct that implements WNJpnRepositoryInterface.
type WNJpnRepository struct {
	SqlProxy sqlproxy.Sql
}

// New is a constructor for WNJpnRepository.
func New(
	sqlProxy sqlproxy.Sql,
) *WNJpnRepository {
	return &WNJpnRepository{
		SqlProxy: sqlProxy,
	}
}

// GetAllAVNWords gets all AVN words.
func (w *WNJpnRepository) GetAllAVNWords(wnJpnDBFilePath string) ([]model.Word, error) {
	return w.getWords(wnJpnDBFilePath, query.GetAllJapaneseAVNWords)
}

// GetAllNWords gets all N words.
func (w *WNJpnRepository) GetAllNWords(wnJpnDBFilePath string) ([]model.Word, error) {
	return w.getWords(wnJpnDBFilePath, query.GetAllJapaneseNWords)
}

// GetAllAVWords gets all AV words.
func (w *WNJpnRepository) GetAllAVWords(wnJpnDBFilePath string) ([]model.Word, error) {
	return w.getWords(wnJpnDBFilePath, query.GetAllJapaneseAVWords)
}

// getWords gets words.
func (w *WNJpnRepository) getWords(wnJpnDBFilePath string, query string) ([]model.Word, error) {
	var deferErr error
	// connect to db
	db, err := w.SqlProxy.Open(sqlproxy.Sqlite, wnJpnDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// execute query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	allWords := make([]model.Word, 0)
	for rows.Next() {
		var word model.Word
		if err := rows.Scan(&word.Lemma, &word.Pos); err != nil {
			return nil, err
		}

		allWords = append(allWords, word)
	}

	return allWords, deferErr
}
