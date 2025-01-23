package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/yanosea/jrp/v2/app/domain/jrp/history"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// HistoryRepository is a struct that implements the HistoryRepository interface.
type historyRepository struct {
	connManager database.ConnectionManager
}

// NewHistoryRepository returns a new instance of the historyRepository struct.
func NewHistoryRepository() history.HistoryRepository {
	return &historyRepository{
		connManager: database.GetConnectionManager(),
	}
}

// DeleteAll is a method that removes all the jrps from the history table.
func (h *historyRepository) DeleteAll(ctx context.Context) (int, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return 0, err
	}

	tx, err := db.BeginTx(
		ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		},
	)
	if err != nil {
		return 0, err
	}
	defer func() {
		deferErr = tx.Rollback()
	}()

	var result proxy.Result
	if result, err = tx.ExecContext(ctx, DeleteAllQuery); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, DeleteSequenceQuery); err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(rowsAffected), deferErr
}

// DeleteByIdIn is a method that removes the jrps from the history table by ID in.
func (h *historyRepository) DeleteByIdIn(ctx context.Context, ids []int) (int, error) {
	var deferErr error
	if len(ids) == 0 {
		return 0, nil
	}

	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return 0, err
	}

	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf(DeleteByIdInQuery, strings.Trim(strings.Repeat("?,", len(ids)), ","))

	var result proxy.Result
	if result, err = db.ExecContext(ctx, query, args...); err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), deferErr
}

// DeleteByIdInAndIsFavoritedIs is a method that removes the jrps from the history table by ID in and is favorited is.
func (h *historyRepository) DeleteByIdInAndIsFavoritedIs(
	ctx context.Context,
	ids []int,
	isFavorited int,
) (int, error) {
	var deferErr error
	if len(ids) == 0 {
		return 0, nil
	}

	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return 0, err
	}

	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(DeleteByIdInAndIsFavoritedIsQuery, strings.Join(placeholders, ","))
	args := make([]interface{}, len(ids)+1)
	for i, id := range ids {
		args[i] = id
	}
	args[len(ids)] = isFavorited

	var result proxy.Result
	if result, err = db.ExecContext(ctx, query, args...); err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), deferErr
}

// DeleteByIsFavoritedIs is a method that removes the jrps from the history table by is favorited.
func (h *historyRepository) DeleteByIsFavoritedIs(ctx context.Context, isFavorited int) (int, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return 0, err
	}

	var result proxy.Result
	if result, err = db.ExecContext(ctx, DeleteByIsFavoritedIsQuery, isFavorited); err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), deferErr
}

// FindAll is a method that finds all the jrps from the history table.
func (h *historyRepository) FindAll(ctx context.Context) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, FindAllQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// FindByIsFavoritedIs is a method that finds the jrps from the history table by is favorited.
func (h *historyRepository) FindByIsFavoritedIs(ctx context.Context, isFavorited int) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, FindByIsFavoritedIsQuery, isFavorited)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// FindByIsFavoritedIsAndPhraseContains is a method that finds the jrps from the history table by is favorited and phrase contains.
func (h *historyRepository) FindByIsFavoritedIsAndPhraseContains(
	ctx context.Context,
	keywords []string,
	and bool,
	isFavorited int,
) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0, len(keywords)+1)
	whereClause := ""
	for i, keyword := range keywords {
		if i == 0 {
			whereClause += "Phrase LIKE ?"
		} else {
			if and {
				whereClause += " AND Phrase LIKE ?"
			} else {
				whereClause += " OR Phrase LIKE ?"
			}
		}
		args = append(args, "%"+keyword+"%")
	}
	args = append(args, isFavorited)
	query := fmt.Sprintf(FindByIsFavoritedIsAndPhraseContainsQuery, whereClause)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// FindByPhraseContains is a method that finds the jrps from the history table by phrase contains.
func (h *historyRepository) FindByPhraseContains(
	ctx context.Context,
	keywords []string,
	and bool,
) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0, len(keywords))
	whereClause := ""
	for i, keyword := range keywords {
		if i == 0 {
			whereClause += "Phrase LIKE ?"
		} else {
			if and {
				whereClause += " AND Phrase LIKE ?"
			} else {
				whereClause += " OR Phrase LIKE ?"
			}
		}
		args = append(args, "%"+keyword+"%")
	}
	query := fmt.Sprintf(FindByPhraseContainsQuery, whereClause)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// FindTopNByIsFavoritedIsAndByOrderByIdAsc is a method that finds the top N jrps from the history table by is favorited order by ID ascending.
func (h *historyRepository) FindTopNByIsFavoritedIsAndByOrderByIdAsc(
	ctx context.Context,
	number int,
	isFavorited int,
) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, FindTopNByIsFavoritedIsAndByOrderByIdAscQuery, isFavorited, number)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc is a method that finds the top N jrps from the history table by is favorited and phrase contains order by ID ascending.
func (h *historyRepository) FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc(
	ctx context.Context,
	keywords []string,
	and bool,
	number int,
	isFavorited int,
) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0, len(keywords)+2)
	whereClause := ""
	for i, keyword := range keywords {
		if i == 0 {
			whereClause += "Phrase LIKE ?"
		} else {
			if and {
				whereClause += " AND Phrase LIKE ?"
			} else {
				whereClause += " OR Phrase LIKE ?"
			}
		}
		args = append(args, "%"+keyword+"%")
	}
	args = append(args, isFavorited, number)
	query := fmt.Sprintf(FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAscQuery, whereClause)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// FindTopNByOrderByIdAsc is a method that finds the top N jrps from the history table by order by ID ascending.
func (h *historyRepository) FindTopNByOrderByIdAsc(ctx context.Context, number int) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, FindTopNByOrderByIdAscQuery, number)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// FindTopNByPhraseContainsOrderByIdAsc is a method that finds the top N jrps from the history table by phrase contains order by ID ascending.
func (h *historyRepository) FindTopNByPhraseContainsOrderByIdAsc(
	ctx context.Context,
	keywords []string,
	and bool,
	number int,
) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0, len(keywords)+1)
	whereClause := ""
	for i, keyword := range keywords {
		if i == 0 {
			whereClause += "Phrase LIKE ?"
		} else {
			if and {
				whereClause += " AND Phrase LIKE ?"
			} else {
				whereClause += " OR Phrase LIKE ?"
			}
		}
		args = append(args, "%"+keyword+"%")
	}
	args = append(args, number)
	query := fmt.Sprintf(FindTopNByPhraseContainsOrderByIdAscQuery, whereClause)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = rows.Close()
	}()

	histories := []*history.History{}
	for rows.Next() {
		history := &history.History{}
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
		histories = append(histories, history)
	}

	return histories, deferErr
}

// SaveAll is a method that saves all the jrp to the history table.
func (h *historyRepository) SaveAll(ctx context.Context, jrps []*history.History) ([]*history.History, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(
		ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = tx.Rollback()
	}()

	stmt, err := db.PrepareContext(ctx, InsertQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr = stmt.Close()
	}()

	for _, jrp := range jrps {
		res, err := stmt.ExecContext(
			ctx,
			jrp.Phrase,
			jrp.Prefix,
			jrp.Suffix,
			jrp.IsFavorited,
			jrp.CreatedAt,
			jrp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		i, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		jrp.ID = int(i)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return jrps, deferErr
}

// UpdateIsFavoritedByIdIn is a method that updates the is favorited of the jrps from the history table by ID in.
func (h *historyRepository) UpdateIsFavoritedByIdIn(
	ctx context.Context,
	isFavorited int,
	ids []int,
) (int, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return 0, err
	}

	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf(UpdateIsFavoritedByIdInQuery, strings.Trim(strings.Repeat("?,", len(ids)), ","))

	var result proxy.Result
	if result, err = db.ExecContext(ctx, query, append([]interface{}{isFavorited}, args...)...); err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), deferErr
}

// UpdateIsFavoritedByIsFavoritedIs is a method that updates the is favorited of the jrps from the history table by is favorited is.
func (h *historyRepository) UpdateIsFavoritedByIsFavoritedIs(
	ctx context.Context,
	isFavorited int,
	isFavoritedIs int,
) (int, error) {
	var deferErr error
	db, err := getJrpDB(ctx, h.connManager)
	if err != nil {
		return 0, err
	}

	var result proxy.Result
	if result, err = db.ExecContext(ctx, UpdateIsFavoritedByIsFavoritedIsQuery, isFavorited, isFavoritedIs); err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), deferErr
}

// getJrpDB is a function that returns the jrp database connection.
func getJrpDB(ctx context.Context, connManager database.ConnectionManager) (proxy.DB, error) {
	var deferErr error
	conn, err := connManager.GetConnection(database.JrpDB)
	if err != nil {
		return nil, err
	}

	db, err := conn.Open()
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, CreateQuery); err != nil {
		return nil, err
	}

	return db, deferErr
}
