package repository

import (
	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/database/jrp/repository/query"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strings"
)

// JrpRepositoryInterface is an interface for JrpRepository.
type JrpRepositoryInterface interface {
	SaveHistory(jrpDBFilePath string, jrps []*model.Jrp) (SaveStatus, error)
	GetAllHistory(jrpDBFilePath string) ([]*model.Jrp, error)
	GetHistoryWithNumber(jrpDBFilePath string, number int) ([]*model.Jrp, error)
	SearchHistoryWithNumber(jrpDBFilePath string, number int, keywords []string, and bool) ([]*model.Jrp, error)
	SearchAllHistory(jrpDBFilePath string, keywords []string, and bool) ([]*model.Jrp, error)
	RemoveHistoryByIDs(jrpDBFilePath string, ids []int, force bool) (RemoveStatus, error)
	RemoveHistoryAll(jrpDBFilePath string, force bool) (RemoveStatus, error)
	GetAllFavorite(jrpDBFilePath string) ([]*model.Jrp, error)
	GetFavoriteWithNumber(jrpDBFilePath string, number int) ([]*model.Jrp, error)
	SearchAllFavorite(jrpDBFilePath string, keywords []string, and bool) ([]*model.Jrp, error)
	SearchFavoriteWithNumber(jrpDBFilePath string, number int, keywords []string, and bool) ([]*model.Jrp, error)
	AddFavoriteByIDs(jrpDBFilePath string, ids []int) (AddStatus, error)
	RemoveFavoriteByIDs(jrpDBFilePath string, ids []int) (RemoveStatus, error)
	RemoveFavoriteAll(jrpDBFilePath string) (RemoveStatus, error)
}

// JrpRepository is a struct that implements JrpRepositoryInterface.
type JrpRepository struct {
	FmtProxy     fmtproxy.Fmt
	SortProxy    sortproxy.Sort
	SqlProxy     sqlproxy.Sql
	StringsProxy stringsproxy.Strings
}

// New is a constructor for JrpRepository.
func New(
	fmtProxy fmtproxy.Fmt,
	sortProxy sortproxy.Sort,
	sqlProxy sqlproxy.Sql,
	stringsProxy stringsproxy.Strings,
) *JrpRepository {
	return &JrpRepository{
		FmtProxy:     fmtProxy,
		SortProxy:    sortProxy,
		SqlProxy:     sqlProxy,
		StringsProxy: stringsProxy,
	}
}

// SaveHistory saves jrps as  history.
func (j JrpRepository) SaveHistory(jrpDBFilePath string, jrps []*model.Jrp) (SaveStatus, error) {
	var deferErr error
	// if jrps is nil or empty, return nil
	if jrps == nil || len(jrps) <= 0 {
		return SavedNone, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return SavedFailed, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return SavedFailed, err
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return SavedFailed, err
	}
	defer func() {
		deferErr = tx.Rollback()
	}()

	// prepare insert statement
	stmt, err := db.Prepare(query.InsertJrp)
	if err != nil {
		return SavedFailed, err
	}
	defer func() {
		deferErr = stmt.Close()
	}()

	// insert jrp and count affected rows
	count := int64(0)
	for _, jrp := range jrps {
		res, err := stmt.Exec(
			jrp.Phrase,
			jrp.Prefix,
			jrp.Suffix,
			jrp.CreatedAt,
			jrp.UpdatedAt,
		)
		if err != nil {
			return SavedFailed, err
		}

		// get count
		c, err := res.RowsAffected()
		if err != nil {
			// failed to get rows affected
			return SavedFailed, err
		}
		count += c

		// set ID
		i, err := res.LastInsertId()
		if err != nil {
			return SavedFailed, err
		}
		jrp.ID = int(i)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return SavedFailed, err
	}

	if count != int64(len(jrps)) {
		// not all rows affected
		return SavedNotAll, nil
	}

	return SavedSuccessfully, deferErr
}

// GetAllHistory gets all jrps as history.
func (j JrpRepository) GetAllHistory(jrpDBFilePath string) ([]*model.Jrp, error) {
	var deferErr error
	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// get all history from jrp
	rows, err := db.Query(query.GetAllJrp)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var allHistory []*model.Jrp
	for rows.Next() {
		history := &model.Jrp{}
		if err := rows.Scan(
			&history.ID,
			&history.Phrase,
			&history.Prefix,
			&history.Suffix,
			&history.IsFavorited,
			&history.CreatedAt,
			&history.UpdatedAt,
		); err != nil {
			return nil, err
		}

		allHistory = append(allHistory, history)
	}

	return allHistory, deferErr
}

// GetHistoryWithNumber gets history with number.
func (j JrpRepository) GetHistoryWithNumber(jrpDBFilePath string, number int) ([]*model.Jrp, error) {
	var deferErr error
	if number <= 0 {
		// if number is less than or equal to 0, return nil
		return nil, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// prepare the query
	stmt, err := db.Prepare(query.GetJrpByNumber)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = stmt.Close()
	}()

	// get history from jrp by number
	rows, err := stmt.Query(number)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var allHistory []*model.Jrp
	for rows.Next() {
		history := &model.Jrp{}
		if err := rows.Scan(
			&history.ID,
			&history.Phrase,
			&history.Prefix,
			&history.Suffix,
			&history.IsFavorited,
			&history.CreatedAt,
			&history.UpdatedAt,
		); err != nil {
			return nil, err
		}

		allHistory = append(allHistory, history)
	}

	// sort by ID asc
	j.SortProxy.Slice(allHistory, func(i, j int) bool {
		return allHistory[i].ID < allHistory[j].ID
	})

	return allHistory, deferErr
}

// SearchAllHistory searches all jrps as history with keywords.
func (j JrpRepository) SearchAllHistory(jrpDBFilePath string, keywords []string, and bool) ([]*model.Jrp, error) {
	var deferErr error
	if keywords == nil || len(keywords) <= 0 {
		// if keywords is nil or empty, return nil
		return nil, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// build query
	args := []interface{}{}
	conditions := []string{}

	// build conditions
	for _, keyword := range keywords {
		conditions = append(conditions, "jrp.Phrase LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	// build where clause
	var whereClause string
	if len(conditions) > 0 {
		separator := " OR "
		if and {
			separator = " AND "
		}
		whereClause = j.StringsProxy.Join(conditions, separator)
	}

	query := j.FmtProxy.Sprintf(query.SearchAllJrp, whereClause)

	// execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var searchedAllHistory []*model.Jrp
	for rows.Next() {
		history := &model.Jrp{}
		if err := rows.Scan(
			&history.ID,
			&history.Phrase,
			&history.Prefix,
			&history.Suffix,
			&history.IsFavorited,
			&history.CreatedAt,
			&history.UpdatedAt,
		); err != nil {
			return nil, err
		}

		searchedAllHistory = append(searchedAllHistory, history)
	}

	return searchedAllHistory, deferErr
}

// SearchHistoryWithNumber searches jrps as history with number and keywords.
func (j JrpRepository) SearchHistoryWithNumber(
	jrpDBFilePath string,
	number int,
	keywords []string,
	and bool,
) ([]*model.Jrp, error) {
	var deferErr error
	if number <= 0 || keywords == nil || len(keywords) <= 0 {
		// if number is less than or equal to 0 or keywords is nil or empty
		return nil, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// build query
	args := []interface{}{}
	conditions := []string{}

	// build conditions
	for _, keyword := range keywords {
		conditions = append(conditions, "jrp.Phrase LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	// build where clause
	var whereClause string
	if len(conditions) > 0 {
		separator := " OR "
		if and {
			separator = " AND "
		}
		whereClause = j.StringsProxy.Join(conditions, separator)
	}

	query := j.FmtProxy.Sprintf(query.SearchJrpByNumber, whereClause)
	args = append(args, number)

	// execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var searchedHistory []*model.Jrp
	for rows.Next() {
		history := &model.Jrp{}
		if err := rows.Scan(
			&history.ID,
			&history.Phrase,
			&history.Prefix,
			&history.Suffix,
			&history.IsFavorited,
			&history.CreatedAt,
			&history.UpdatedAt,
		); err != nil {
			return nil, err
		}

		searchedHistory = append(searchedHistory, history)
	}

	// sort by ID asc
	j.SortProxy.Slice(searchedHistory, func(i, j int) bool {
		return searchedHistory[i].ID < searchedHistory[j].ID
	})

	return searchedHistory, deferErr
}

// RemoveHistoryByIDs removes jrps by IDs.
func (j JrpRepository) RemoveHistoryByIDs(jrpDBFilePath string, ids []int, force bool) (RemoveStatus, error) {
	var deferErr error
	if ids == nil || len(ids) <= 0 {
		// if ids is nil or empty, return nil
		return RemovedNone, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return RemovedFailed, err
	}

	// create the correct number of placeholders for the IDs
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	placeholdersStr := j.StringsProxy.Join(placeholders, ",")

	// prepare the delete q with the correct number of placeholders
	var q string
	if force {
		q = j.FmtProxy.Sprintf(query.RemoveJrpByIDs, placeholdersStr)
	} else {
		q = j.FmtProxy.Sprintf(query.RemoveJrpByIDsExceptFavorite, placeholdersStr)
	}
	stmt, err := db.Prepare(q)
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = stmt.Close()
	}()

	// convert ids to interface slice for Exec
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// execute the delete query
	res, err := stmt.Exec(args...)
	if err != nil {
		return RemovedFailed, err
	}

	if count, err := res.RowsAffected(); err != nil {
		// failed to get rows affected
		return RemovedFailed, err
	} else if count <= 0 {
		// no rows affected
		return RemovedNone, nil
	} else if count != int64(len(ids)) {
		// not all rows affected
		return RemovedNotAll, nil
	}

	return RemovedSuccessfully, deferErr
}

// RemoveHistoryAll removes all jrps.
func (j JrpRepository) RemoveHistoryAll(jrpDBFilePath string, force bool) (RemoveStatus, error) {
	var deferErr error
	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return RemovedFailed, err
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = tx.Rollback()
	}()

	// set q
	var q string
	if force {
		q = query.RemoveAllJrp
	} else {
		q = query.RemoveAllJrpExceptFavorite
	}

	// remove all jrp
	var res sqlproxy.ResultInstanceInterface
	if res, err = tx.Exec(q); err != nil {
		return RemovedFailed, err
	}

	// check if the count of rows is zero after execution
	checkCount := query.CountJrp
	var count int
	if err := tx.QueryRow(checkCount).Scan(&count); err != nil {
		return RemovedFailed, err
	}

	if count == 0 {
		// remove jrp sequence
		if _, err := tx.Exec(query.RemoveJrpSeq); err != nil {
			return RemovedFailed, err
		}
	}

	// check rows affected by the remove jrp query
	if affected, err := res.RowsAffected(); err != nil {
		return RemovedFailed, err
	} else if affected == 0 {
		return RemovedNone, nil
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return RemovedFailed, err
	}

	return RemovedSuccessfully, deferErr
}

// GetAllFavorite gets all jrps that are favorited.
func (j JrpRepository) GetAllFavorite(jrpDBFilePath string) ([]*model.Jrp, error) {
	var deferErr error
	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// get all favorite from jrp
	rows, err := db.Query(query.GetAllFavorite)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var allFavorite []*model.Jrp
	for rows.Next() {
		favorite := &model.Jrp{}
		if err := rows.Scan(
			&favorite.ID,
			&favorite.Phrase,
			&favorite.Prefix,
			&favorite.Suffix,
			&favorite.IsFavorited,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
		); err != nil {
			return nil, err
		}

		allFavorite = append(allFavorite, favorite)
	}

	return allFavorite, deferErr
}

// GetFavoriteWithNumber gets jrps that are favorited with number.
func (j JrpRepository) GetFavoriteWithNumber(jrpDBFilePath string, number int) ([]*model.Jrp, error) {
	var deferErr error
	if number <= 0 {
		// if number is less than or equal to 0, return nil
		return nil, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// prepare the query
	stmt, err := db.Prepare(query.GetFavoriteByNumber)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = stmt.Close()
	}()

	// get favorite from jrp by number
	rows, err := stmt.Query(number)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var allFavorite []*model.Jrp
	for rows.Next() {
		favorite := &model.Jrp{}
		if err := rows.Scan(
			&favorite.ID,
			&favorite.Phrase,
			&favorite.Prefix,
			&favorite.Suffix,
			&favorite.IsFavorited,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
		); err != nil {
			return nil, err
		}

		allFavorite = append(allFavorite, favorite)
	}

	// sort by ID asc
	j.SortProxy.Slice(allFavorite, func(i, j int) bool {
		return allFavorite[i].ID < allFavorite[j].ID
	})

	return allFavorite, deferErr
}

// SearchAllFavorite searches all jrps that are favorited with keywords.
func (j JrpRepository) SearchAllFavorite(jrpDBFilePath string, keywords []string, and bool) ([]*model.Jrp, error) {
	var deferErr error
	if keywords == nil || len(keywords) <= 0 {
		// if keywords is nil or empty, return nil
		return nil, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// build query
	args := []interface{}{}
	conditions := []string{}

	// build conditions
	for _, keyword := range keywords {
		conditions = append(conditions, "jrp.Phrase LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	// build where clause
	var whereClause string
	if len(conditions) > 0 {
		separator := " OR "
		if and {
			separator = " AND "
		}
		whereClause = j.StringsProxy.Join(conditions, separator)
	}

	query := j.FmtProxy.Sprintf(query.SearchAllFavorite, whereClause)

	// execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var searchedAllFavorite []*model.Jrp
	for rows.Next() {
		favorite := &model.Jrp{}
		if err := rows.Scan(&favorite.ID,
			&favorite.Phrase,
			&favorite.Prefix,
			&favorite.Suffix,
			&favorite.IsFavorited,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
		); err != nil {
			return nil, err
		}

		searchedAllFavorite = append(searchedAllFavorite, favorite)
	}

	return searchedAllFavorite, deferErr
}

// SearchFavoriteWithNumber searches jrps that are favorited with number and keywords.
func (j JrpRepository) SearchFavoriteWithNumber(
	jrpDBFilePath string,
	number int,
	keywords []string,
	and bool,
) ([]*model.Jrp, error) {
	var deferErr error
	if number <= 0 || keywords == nil || len(keywords) <= 0 {
		// if number is less than or equal to 0 or keywords is nil or empty
		return nil, nil
	}

	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return nil, err
	}

	// build query
	args := []interface{}{}
	conditions := []string{}

	// build conditions
	for _, keyword := range keywords {
		conditions = append(conditions, "jrp.Phrase LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	// build where clause
	var whereClause string
	if len(conditions) > 0 {
		separator := " OR "
		if and {
			separator = " AND "
		}
		whereClause = j.StringsProxy.Join(conditions, separator)
	}

	query := j.FmtProxy.Sprintf(query.SearchFavoriteByNumber, whereClause)
	args = append(args, number)

	// execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	// scan rows
	var searchedFavorite []*model.Jrp
	for rows.Next() {
		favorite := &model.Jrp{}
		if err := rows.Scan(
			&favorite.ID,
			&favorite.Phrase,
			&favorite.Prefix,
			&favorite.Suffix,
			&favorite.IsFavorited,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
		); err != nil {
			return nil, err
		}

		searchedFavorite = append(searchedFavorite, favorite)
	}

	// sort by ID asc
	j.SortProxy.Slice(searchedFavorite, func(i, j int) bool {
		return searchedFavorite[i].ID < searchedFavorite[j].ID
	})

	return searchedFavorite, deferErr
}

// AddFavoriteByIDs adds jrps to favorite by IDs.
func (j JrpRepository) AddFavoriteByIDs(jrpDBFilePath string, ids []int) (AddStatus, error) {
	var deferErr error
	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return AddedFailed, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return AddedFailed, err
	}

	// create the correct number of placeholders for the IDs
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	placeholdersStr := j.StringsProxy.Join(placeholders, ",")

	// prepare the delete query with the correct number of placeholders
	query := j.FmtProxy.Sprintf(query.AddFavoriteByIDs, placeholdersStr)
	stmt, err := db.Prepare(query)
	if err != nil {
		return AddedFailed, err
	}

	// convert ids to interface slice for Exec
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// execute the delete query
	res, err := stmt.Exec(args...)
	if err != nil {
		return AddedFailed, err
	}

	if count, err := res.RowsAffected(); err != nil {
		// failed to get rows affected
		return AddedFailed, err
	} else if count <= 0 {
		// no rows affected
		return AddedNone, nil
	} else if count != int64(len(ids)) {
		// not all rows affected
		return AddedNotAll, nil
	}

	return AddedSuccessfully, deferErr
}

// RemoveFavoriteByIDs removes jrps from favorite by IDs.
func (j JrpRepository) RemoveFavoriteByIDs(jrpDBFilePath string, ids []int) (RemoveStatus, error) {
	var deferErr error
	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return RemovedFailed, err
	}

	// create the correct number of placeholders for the IDs
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	placeholdersStr := j.StringsProxy.Join(placeholders, ",")

	// prepare the delete query with the correct number of placeholders
	query := j.FmtProxy.Sprintf(query.RemoveFavoriteByIDs, placeholdersStr)
	stmt, err := db.Prepare(query)
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = stmt.Close()
	}()

	// convert ids to interface slice for Exec
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// execute the delete query
	res, err := stmt.Exec(args...)
	if err != nil {
		return RemovedFailed, err
	}

	// check rows affected
	if count, err := res.RowsAffected(); err != nil {
		// failed to get rows affected
		return RemovedFailed, err
	} else if count <= 0 {
		// no rows affected
		return RemovedNone, nil
	} else if count != int64(len(ids)) {
		// not all rows affected
		return RemovedNotAll, nil
	}

	return RemovedSuccessfully, deferErr
}

// RemoveFavoriteAll removes all jrps from favorite.
func (j JrpRepository) RemoveFavoriteAll(jrpDBFilePath string) (RemoveStatus, error) {
	var deferErr error
	// connect to db
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = db.Close()
	}()

	// create table 'jrp'
	if _, err := j.createTableJrp(db); err != nil {
		return RemovedFailed, err
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return RemovedFailed, err
	}
	defer func() {
		deferErr = tx.Rollback()
	}()

	// remove all favorite
	var res sqlproxy.ResultInstanceInterface
	if res, err = tx.Exec(query.RemoveAllFavorite); err != nil {
		return RemovedFailed, err
	}

	// check rows affected
	if count, err := res.RowsAffected(); err != nil {
		// failed to get rows affected
		return RemovedFailed, err
	} else if count <= 0 {
		// no rows affected
		return RemovedNone, err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return RemovedFailed, err
	}

	return RemovedSuccessfully, deferErr
}

// createTableJrp creates table 'jrp'.
func (j JrpRepository) createTableJrp(db sqlproxy.DBInstanceInterface) (sqlproxy.ResultInstanceInterface, error) {
	return db.Exec(query.CreateTableJrp)
}
