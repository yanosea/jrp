package jrp

import (
	"context"
	"time"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"
)

// getHistoryUseCase is a struct that contains the use case of the getting jrp from the table history in jrp sqlite database.
type getHistoryUseCase struct {
	historyRepo historyDomain.HistoryRepository
}

// NewGetHistoryUseCase returns a new instance of the GetHistoryUseCase struct.
func NewGetHistoryUseCase(
	historyRepo historyDomain.HistoryRepository,
) *getHistoryUseCase {
	return &getHistoryUseCase{
		historyRepo: historyRepo,
	}
}

// GetHistoryUseCaseOutputDto is a DTO struct that contains the output data of the GetHistoryUseCase.
type GetHistoryUseCaseOutputDto struct {
	// ID is the identifier of the phrase
	ID int
	// Phrase is the generated phrase
	Phrase string
	// Prefix is the prefix when the phrase is generated
	Prefix string
	// Suffix is the suffix when the phrase is generated
	Suffix string
	// IsFavorited is the flag to indicate whether the phrase is favorited
	IsFavorited int
	// CreatedAt is the timestamp when the phrase is created
	CreatedAt time.Time
	// UpdatedAt is the timestamp when the phrase is updated
	UpdatedAt time.Time
}

// Run returns the output of the GetHistoryUseCase.
func (uc *getHistoryUseCase) Run(ctx context.Context, all bool, favorited bool, number int) ([]*GetHistoryUseCaseOutputDto, error) {
	var histories []*historyDomain.History
	var err error
	if all && favorited {
		histories, err = uc.historyRepo.FindByIsFavoritedIs(ctx, 1)
	} else if all && !favorited {
		histories, err = uc.historyRepo.FindAll(ctx)
	} else if !all && favorited {
		histories, err = uc.historyRepo.FindTopNByIsFavoritedIsAndByOrderByIdAsc(ctx, number, 1)
	} else {
		histories, err = uc.historyRepo.FindTopNByOrderByIdAsc(ctx, number)
	}
	if err != nil {
		return nil, err
	}

	var ucDtos []*GetHistoryUseCaseOutputDto
	for _, history := range histories {
		ucDtos = append(ucDtos, &GetHistoryUseCaseOutputDto{
			ID:          history.ID,
			Phrase:      history.Phrase,
			Prefix:      history.Prefix.String,
			Suffix:      history.Suffix.String,
			IsFavorited: history.IsFavorited,
			CreatedAt:   history.CreatedAt,
			UpdatedAt:   history.UpdatedAt,
		})
	}

	return ucDtos, nil
}
