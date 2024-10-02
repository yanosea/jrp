package query

import (
	_ "embed"
)

// CreateTableJrp is a query to create table jrp.
//
//go:embed create_table_jrp.sql
var CreateTableJrp string

// InsertJrp is a query to insert jrp.
//
//go:embed insert_jrp.sql
var InsertJrp string

// GetAllJrp is a query to get all jrp.
//
//go:embed get_all_jrp.sql
var GetAllJrp string

// GetJrpByNumber is a query to get jrp by number.
//
//go:embed get_jrp_by_number.sql
var GetJrpByNumber string

// RemoveJrpByIDs is a query to remove jrp by ids.
//
//go:embed remove_jrp_by_ids.sql
var RemoveJrpByIDs string

// RemoveJrpByIDsExceptFavorite is a query to remove jrp by ids except favorite.
//
//go:embed remove_jrp_by_ids_except_favorite.sql
var RemoveJrpByIDsExceptFavorite string

// RemoveAllJrp is a query to remove all jrp.
//
//go:embed remove_all_jrp.sql
var RemoveAllJrp string

// RemoveAllJrpExceptFavorite is a query to remove all jrp except favorite.
//
//go:embed remove_all_jrp_except_favorite.sql
var RemoveAllJrpExceptFavorite string

// CountJrp is a query to count jrp.
//
//go:embed count_jrp.sql
var CountJrp string

// RemoveJrpSeq is a query to remove jrp seq.
//
//go:embed remove_jrp_seq.sql
var RemoveJrpSeq string

// SearchAllJrp is a query to search all jrp.
//
//go:embed search_all_jrp.sql
var SearchAllJrp string

// SearchJrpByNumber is a query to search jrp by number.
//
//go:embed search_jrp_by_number.sql
var SearchJrpByNumber string

// GetAllFavorite is a query to get all favorite.
//
//go:embed get_all_favorite.sql
var GetAllFavorite string

// GetFavoriteByNumber is a query to get favorite by number.
//
//go:embed get_favorite_by_number.sql
var GetFavoriteByNumber string

// AddFavoriteByIDs is a query to add favorite by ids.
//
//go:embed add_favorite_by_ids.sql
var AddFavoriteByIDs string

// RemoveFavoriteByIDs is a query to remove favorite by ids.
//
//go:embed remove_favorite_by_ids.sql
var RemoveFavoriteByIDs string

// RemoveAllFavorite is a query to remove all favorite.
//
//go:embed remove_all_favorite.sql
var RemoveAllFavorite string

// SearchAllFavorite is a query to search all favorite.
//
//go:embed search_all_favorite.sql
var SearchAllFavorite string

// SearchFavoriteByNumber is a query to search favorite by number.
//
//go:embed search_favorite_by_number.sql
var SearchFavoriteByNumber string
