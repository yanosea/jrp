package jrp

import (
	"context"
	"errors"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"
)

// favoriteUseCase is a struct that contains the use case of the favoriting jrp from the table history in jrp sqlite database.
type favoriteUseCase struct {
	historyRepo historyDomain.HistoryRepository
}

// NewFavoriteUseCase returns a new instance of the FavoriteUseCase struct.
func NewFavoriteUseCase(
	historyRepo historyDomain.HistoryRepository,
) *favoriteUseCase {
	return &favoriteUseCase{
		historyRepo: historyRepo,
	}
}

// Run returns the output of the AddFavorite usecase.
func (uc *favoriteUseCase) Run(ctx context.Context, ids []int, all bool) error {
	var rowsAffected int
	var err error
	if all {
		rowsAffected, err = uc.historyRepo.UpdateIsFavoritedByIsFavoritedIs(ctx, 1, 0)
	} else {
		rowsAffected, err = uc.historyRepo.UpdateIsFavoritedByIdIn(ctx, 1, ids)
	}
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no histories to favorite")
	}

	return nil
}
