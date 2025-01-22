package query_service

import ()

const (
	// FindByLangIsAndPosInQuery is a query that finds the records from the word table by lang is and pos in.
	FindByLangIsAndPosInQuery = `
SELECT
    word.WordID
    , word.Lang
    , word.Lemma
    , word.Pron
    , word.Pos
FROM
    word
WHERE
    word.Lang = ?
    AND word.Pos IN (%s);
`
)
