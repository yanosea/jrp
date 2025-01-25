package history

import (
	"database/sql"
	"time"
)

// History is a struct that represents history table in the jrp database.
type History struct {
	// ID is the primary key of the history table.
	ID int
	// Phrase is the generated phrase.
	Phrase string
	// Prefix is the prefix when the phrase is generated.
	Prefix sql.NullString
	// Suffix is the suffix when the phrase is generated.
	Suffix sql.NullString
	// IsFavorited is the flag to indicate whether the phrase is favorited.
	IsFavorited int
	// CreatedAt is the timestamp when the phrase is created.
	CreatedAt time.Time
	// UpdatedAt is the timestamp when the phrase is updated.
	UpdatedAt time.Time
}

// NewHistory returns a new instance of the History struct.
func NewHistory(
	phrase string,
	prefix string,
	suffix string,
	isFavorited int,
	createdAt time.Time,
	updatedAt time.Time,
) *History {
	return &History{
		Phrase:      phrase,
		Prefix:      sql.NullString{String: prefix, Valid: prefix != ""},
		Suffix:      sql.NullString{String: suffix, Valid: suffix != ""},
		IsFavorited: isFavorited,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
