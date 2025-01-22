package repository

import ()

const (
	// CreateQuery is a query that creates a table history.
	CreateQuery = `
CREATE TABLE IF NOT EXISTS
  history (
    ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT
    , Phrase TEXT NOT NULL
    , Prefix TEXT
    , Suffix TEXT
    , IsFavorited INTEGER DEFAULT 0
    , CreatedAt TIMESTAMP
    , UpdatedAt TIMESTAMP
  );
`
	// DeleteAllQuery is a query that deletes all from the history table.
	DeleteAllQuery = `
DELETE
FROM
  history;
`
	// DeleteByIdInQuery is a query that deletes the records from the history table by ID in.
	DeleteByIdInQuery = `
DELETE
FROM
  history
WHERE
  history.ID IN (%s);
`
	// DeleteByIdInAndIsFavoritedIsQuery is a query that deletes the records from the history table by ID in and is favorited is.
	DeleteByIdInAndIsFavoritedIsQuery = `
DELETE
FROM
  history
WHERE
  history.ID IN (%s)
  AND history.IsFavorited = ?;
`
	// DeleteByIdInQuery is a query that deletes the records from the history table by ID in.
	DeleteByIsFavoritedIsQuery = `
DELETE
FROM
  history
WHERE
  history.IsFavorited = ?;
`
	// DeleteByIsFavoritedIsQuery is a query that deletes the records from the history table by is favorited.
	DeleteSequenceQuery = `
DELETE
FROM
  sqlite_sequence
WHERE
  sqlite_sequence.name = 'history';
`
	// FindAllQuery is a query that finds all from the history table.
	FindAllQuery = `
SELECT
  history.ID
  , history.Phrase
  , history.Prefix
  , history.Suffix
  , history.IsFavorited
  , history.CreatedAt
  , history.UpdatedAt
FROM
  history
ORDER BY
  history.ID ASC;
`
	// FindByIsFavoritedIsQuery is a query that finds the records from the history table by is favorited.
	FindByIsFavoritedIsQuery = `
SELECT
  history.ID
  , history.Phrase
  , history.Prefix
  , history.Suffix
  , history.IsFavorited
  , history.CreatedAt
  , history.UpdatedAt
FROM
  history
WHERE
  history.IsFavorited = ?
ORDER BY
  history.ID ASC;
`
	// FindByIsFavoritedIsAndPhraseContainsQuery is a query that finds the records from the history table by is favorited and phrase contains.
	FindByIsFavoritedIsAndPhraseContainsQuery = `
SELECT
  history.ID
  , history.Phrase
  , history.Prefix
  , history.Suffix
  , history.IsFavorited
  , history.CreatedAt
  , history.UpdatedAt
FROM
  history
WHERE
  (%s)
  AND history.IsFavorited = ?
ORDER BY
  history.ID ASC;
`
	// FindByPhraseContainsQuery is a query that finds the records from the history table by phrase contains.
	FindByPhraseContainsQuery = `
SELECT
  history.ID
  , history.Phrase
  , history.Prefix
  , history.Suffix
  , history.IsFavorited
  , history.CreatedAt
  , history.UpdatedAt
FROM
  history
WHERE
  (%s)
ORDER BY
  history.ID ASC;
`
	// FindTopNByIsFavoritedIsAndByOrderByIdAscQuery is a query that finds the top N records from the history table by is favorited order by ID ascending.
	FindTopNByIsFavoritedIsAndByOrderByIdAscQuery = `
SELECT
  *
FROM (
  SELECT
    history.ID
    , history.Phrase
    , history.Prefix
    , history.Suffix
    , history.IsFavorited
    , history.CreatedAt
    , history.UpdatedAt
  FROM
    history
  WHERE
    history.IsFavorited = ?
  ORDER BY
    history.ID DESC
  LIMIT ?
) AS latest_records
ORDER BY
  latest_records.ID ASC;
`
	// FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAscQuery is a query that finds the top N records from the history table by is favorited and phrase contains order by ID ascending.
	FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAscQuery = `
SELECT
  *
FROM (
  SELECT
    history.ID
    , history.Phrase
    , history.Prefix
    , history.Suffix
    , history.IsFavorited
    , history.CreatedAt
    , history.UpdatedAt
  FROM
    history
  WHERE
    (%s)
    AND history.IsFavorited = ?
  ORDER BY
    history.ID DESC
  LIMIT ?
) AS latest_records
ORDER BY
  latest_records.ID ASC;
`
	// FindTopNByOrderByIdAscQuery is a query that finds the top N records from the history table by order by ID ascending.
	FindTopNByOrderByIdAscQuery = `
SELECT
  *
FROM (
  SELECT
    history.ID
    , history.Phrase
    , history.Prefix
    , history.Suffix
    , history.IsFavorited
    , history.CreatedAt
    , history.UpdatedAt
  FROM
    history
  ORDER BY
    history.ID DESC
  LIMIT ?
) AS latest_records
ORDER BY
  latest_records.ID ASC;
`
	// FindTopNByPhraseContainsOrderByIdAscQuery is a query that finds the top N records from the history table by phrase contains order by ID ascending.
	FindTopNByPhraseContainsOrderByIdAscQuery = `
SELECT
  *
FROM (
  SELECT
    history.ID
    , history.Phrase
    , history.Prefix
    , history.Suffix
    , history.IsFavorited
    , history.CreatedAt
    , history.UpdatedAt
  FROM
    history
  WHERE
    (%s)
  ORDER BY
    history.ID DESC
  LIMIT ?
) AS latest_records
ORDER BY
  latest_records.ID ASC;
`
	// InsertQuery is a query that inserts a record into the history table.
	InsertQuery = `
INSERT INTO
  history (
    Phrase
    , Prefix
    , Suffix
    , IsFavorited
    , CreatedAt
    , UpdatedAt
  ) VALUES (
    ?
    , ?
    , ?
    , ?
    , ?
    , ?
  );
`
	// UpdateIsFavoritedByIdInQuery is a query that updates the is favorited by ID in.
	UpdateIsFavoritedByIdInQuery = `
UPDATE
  history
SET
  IsFavorited = ?
WHERE
  history.ID IN (%s);
`
	// UpdateIsFavoritedByIsFavoritedIsQuery is a query that updates the is favorited by is favorited.
	UpdateIsFavoritedByIsFavoritedIsQuery = `
UPDATE
  history
SET
  IsFavorited = ?
WHERE
  history.IsFavorited = ?;
`
)
