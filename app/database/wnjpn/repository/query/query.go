package query

import (
	_ "embed"
)

// GetAllJapaneseAVNWords is a query to get all japanese avn words.
//
//go:embed get_all_avn_words.sql
var GetAllJapaneseAVNWords string

// GetAllJapaneseNWords is a query to get all japanese n words.
//
//go:embed get_all_n_words.sql
var GetAllJapaneseNWords string

// GetAllJapaneseAVWords is a query to get all japanese av words.
//
//go:embed get_all_av_words.sql
var GetAllJapaneseAVWords string
