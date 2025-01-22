package jrp

import (
	"context"
	"errors"

	historyDomain "github.com/yanosea/jrp/app/domain/jrp/history"
)

// removeHistoryUseCase is a struct that contains the use case of the removing jrp from the table history in jrp sqlite database.
type removeHistoryUseCase struct {
	historyRepo historyDomain.HistoryRepository
}

// NewRemoveHistoryUseCase returns a new instance of the RemoveHistoryUseCase struct.
func NewRemoveHistoryUseCase(
	historyRepo historyDomain.HistoryRepository,
) *removeHistoryUseCase {
	return &removeHistoryUseCase{
		historyRepo: historyRepo,
	}
}

// Run returns the output of the RemoveHistoryUseCase.
func (uc *removeHistoryUseCase) Run(ctx context.Context, ids []int, all bool, force bool) error {
	var rowsAffected int
	var err error
	if all && force {
		rowsAffected, err = uc.historyRepo.DeleteAll(ctx)
	} else if all && !force {
		rowsAffected, err = uc.historyRepo.DeleteByIsFavoritedIs(ctx, 0)
	} else if !all && force {
		rowsAffected, err = uc.historyRepo.DeleteByIdIn(ctx, ids)
	} else if !all && !force {
		rowsAffected, err = uc.historyRepo.DeleteByIdInAndIsFavoritedIs(ctx, ids, 0)
	}
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no histories to remove")
	}

	return nil
}
