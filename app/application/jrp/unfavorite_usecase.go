package jrp

import (
	"context"
	"errors"

	historyDomain "github.com/yanosea/jrp/app/domain/jrp/history"
)

// unfavoriteUseCase is a struct that contains the use case of the unfavoriting jrp from the table history in jrp sqlite database.
type unfavoriteUseCase struct {
	historyRepo historyDomain.HistoryRepository
}

// NewUnfavoriteUseCase returns a new instance of the UnfavoriteUseCase struct.
func NewUnfavoriteUseCase(
	historyRepo historyDomain.HistoryRepository,
) *unfavoriteUseCase {
	return &unfavoriteUseCase{
		historyRepo: historyRepo,
	}
}

// Run returns the output of the Unfavorite usecase.
func (uc *unfavoriteUseCase) Run(ctx context.Context, ids []int, all bool) error {
	var rowsAffected int
	var err error
	if all {
		rowsAffected, err = uc.historyRepo.UpdateIsFavoritedByIsFavoritedIs(ctx, 0, 1)
	} else {
		rowsAffected, err = uc.historyRepo.UpdateIsFavoritedByIdIn(ctx, 0, ids)
	}
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no favorited histories to unfavorite")
	}

	return nil
}
