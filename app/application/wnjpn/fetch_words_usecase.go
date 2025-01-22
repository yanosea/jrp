package wnjpn

import (
	"context"
)

// FetchWordsUseCase is an interface that defines the use case of fetching words.
type FetchWordsUseCase interface {
	Run(ctx context.Context, lang string, pos []string) ([]*FetchWordsUseCaseOutputDto, error)
}

// FetchWordsUseCaseStruct is a struct that implements the FetchWordsUseCase interface.
type FetchWordsUseCaseStruct struct {
	wordQueryService WordQueryService
}

var (
	// NewFetchWordsUseCase is a function that returns a new instance of the fetchWordsUseCase struct.
	NewFetchWordsUseCase = newFetchWordsUseCase
)

// newFetchWordsUseCase returns a new instance of the fetchWordsUseCase struct.
func newFetchWordsUseCase(
	wordQueryService WordQueryService,
) *FetchWordsUseCaseStruct {
	return &FetchWordsUseCaseStruct{
		wordQueryService: wordQueryService,
	}
}

// FetchWordsUseCaseOutputDto is a DTO struct that contains the output data of the FetchWordsUseCase.
type FetchWordsUseCaseOutputDto struct {
	WordID int
	Lang   string
	Lemma  string
	Pron   string
	Pos    string
}

// Run returns the output of the FetchWordsUseCase.
func (uc *FetchWordsUseCaseStruct) Run(ctx context.Context, lang string, pos []string) ([]*FetchWordsUseCaseOutputDto, error) {
	qsDtos, err := uc.wordQueryService.FindByLangIsAndPosIn(ctx, lang, pos)
	if err != nil {
		return nil, err
	}

	var ucDtos []*FetchWordsUseCaseOutputDto
	for _, qsDto := range qsDtos {
		ucDtos = append(ucDtos, &FetchWordsUseCaseOutputDto{
			WordID: qsDto.WordID,
			Lang:   qsDto.Lang.String,
			Lemma:  qsDto.Lemma.String,
			Pron:   qsDto.Pron.String,
			Pos:    qsDto.Pos.String,
		})
	}

	return ucDtos, nil
}
