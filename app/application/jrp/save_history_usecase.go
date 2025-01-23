package jrp

import (
	"context"
	"time"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"
)

// saveHistoryUseCase is a struct that contains the use case of the saving jrp to the table history in jrp sqlite database.
type saveHistoryUseCase struct {
	historyRepo historyDomain.HistoryRepository
}

// NewSaveHistoryUseCase returns a new instance of the SaveHistoryUseCase struct.
func NewSaveHistoryUseCase(
	historyRepo historyDomain.HistoryRepository,
) *saveHistoryUseCase {
	return &saveHistoryUseCase{
		historyRepo: historyRepo,
	}
}

// SaveHistoryUseCaseInputDto is a DTO struct that contains the output data of the SaveHistoryUseCase.
type SaveHistoryUseCaseInputDto struct {
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

// SaveHistoryUseCaseOutputDto is a DTO struct that contains the output data of the SaveHistoryUseCase.
type SaveHistoryUseCaseOutputDto struct {
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

// Run returns the output of the SaveHistoryUseCase.
func (uc *saveHistoryUseCase) Run(ctx context.Context, inputDtos []*SaveHistoryUseCaseInputDto) ([]*SaveHistoryUseCaseOutputDto, error) {
	var histories []*historyDomain.History
	for _, dto := range inputDtos {
		history := historyDomain.NewHistory(
			dto.Phrase,
			dto.Prefix,
			dto.Suffix,
			dto.IsFavorited,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		histories = append(histories, history)
	}

	histories, err := uc.historyRepo.SaveAll(ctx, histories)
	if err != nil {
		return nil, err
	}

	var outputDtos []*SaveHistoryUseCaseOutputDto
	for _, history := range histories {
		outputDto := &SaveHistoryUseCaseOutputDto{
			ID:          history.ID,
			Phrase:      history.Phrase,
			Prefix:      history.Prefix.String,
			Suffix:      history.Suffix.String,
			IsFavorited: history.IsFavorited,
			CreatedAt:   history.CreatedAt,
			UpdatedAt:   history.UpdatedAt,
		}
		outputDtos = append(outputDtos, outputDto)
	}

	return outputDtos, nil
}
