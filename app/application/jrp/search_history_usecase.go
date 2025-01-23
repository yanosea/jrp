package jrp

import (
	"context"
	"time"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"
)

// searchHistoryUseCase is a struct that contains the use case of the searching jrp from the table history in jrp sqlite database.
type searchHistoryUseCase struct {
	historyRepo historyDomain.HistoryRepository
}

// NewSearchHistoryUseCase returns a new instance of the SearchHistoryUseCase struct.
func NewSearchHistoryUseCase(
	historyRepo historyDomain.HistoryRepository,
) *searchHistoryUseCase {
	return &searchHistoryUseCase{
		historyRepo: historyRepo,
	}
}

// SearchHistoryUseCaseOutputDto is a DTO struct that contains the output data of the SearchHistoryUseCase.
type SearchHistoryUseCaseOutputDto struct {
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

// Run returns the output of the SearchHistoryUseCase.
func (uc *searchHistoryUseCase) Run(ctx context.Context, keywords []string, and bool, all bool, favorited bool, number int) ([]*SearchHistoryUseCaseOutputDto, error) {
	var histories []*historyDomain.History
	var err error
	if all && favorited {
		histories, err = uc.historyRepo.FindByIsFavoritedIsAndPhraseContains(ctx, keywords, and, 1)
	} else if all && !favorited {
		histories, err = uc.historyRepo.FindByPhraseContains(ctx, keywords, and)
	} else if !all && favorited {
		histories, err = uc.historyRepo.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc(ctx, keywords, and, number, 1)
	} else {
		histories, err = uc.historyRepo.FindTopNByPhraseContainsOrderByIdAsc(ctx, keywords, and, number)
	}
	if err != nil {
		return nil, err
	}

	var ucDtos []*SearchHistoryUseCaseOutputDto
	for _, h := range histories {
		ucDtos = append(ucDtos, &SearchHistoryUseCaseOutputDto{
			ID:          h.ID,
			Phrase:      h.Phrase,
			Prefix:      h.Prefix.String,
			Suffix:      h.Suffix.String,
			IsFavorited: h.IsFavorited,
			CreatedAt:   h.CreatedAt,
			UpdatedAt:   h.UpdatedAt,
		})
	}

	return ucDtos, nil
}
