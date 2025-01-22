package jrp

import (
	"time"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"
)

// generateJrpUseCase is a struct that contains the use case of the generation jrp.
type generateJrpUseCase struct{}

// NewGenerateJrpUseCase returns a new instance of the GenerateJrpUseCase struct.
func NewGenerateJrpUseCase() *generateJrpUseCase {
	return &generateJrpUseCase{}
}

// GenerateJrpUseCaseInputDto is a DTO struct that contains the input data of the GenerateJrpUseCase.
type GenerateJrpUseCaseInputDto struct {
	WordID int
	Lang   string
	Lemma  string
	Pron   string
	Pos    string
}

// GenerateJrpUseCaseOutputDto is a DTO struct that contains the output data of the GenerateJrpUseCase.
type GenerateJrpUseCaseOutputDto struct {
	ID          int
	Phrase      string
	Prefix      string
	Suffix      string
	IsFavorited int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	// ru is a variable that contains the RandUtil struct for injecting dependencies in testing.
	ru = utility.NewRandUtil(proxy.NewRand())
)

// RunWithPrefix generates a jrp with the given prefix.
func (uc *generateJrpUseCase) RunWithPrefix(
	dtos []*GenerateJrpUseCaseInputDto,
	prefix string,
) *GenerateJrpUseCaseOutputDto {
	if len(dtos) == 0 {
		return nil
	}

	now := time.Now()
	maxAttempts := len(dtos)

	var jrp *GenerateJrpUseCaseOutputDto = nil
	for i := 0; i < maxAttempts; i++ {
		randomSuffix := dtos[ru.GenerateRandomNumber(len(dtos))]
		if randomSuffix.Pos != "n" {
			continue
		}

		jrp = &GenerateJrpUseCaseOutputDto{
			ID:          0,
			Phrase:      prefix + randomSuffix.Lemma,
			Prefix:      prefix,
			Suffix:      "",
			IsFavorited: 0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		break
	}

	return jrp
}

// RunWithSuffix generates a jrp with the given suffix.
func (uc *generateJrpUseCase) RunWithSuffix(
	dtos []*GenerateJrpUseCaseInputDto,
	suffix string,
) *GenerateJrpUseCaseOutputDto {
	if len(dtos) == 0 {
		return nil
	}

	now := time.Now()
	maxAttempts := len(dtos)

	var jrp *GenerateJrpUseCaseOutputDto = nil
	for i := 0; i < maxAttempts; i++ {
		randomPrefix := dtos[ru.GenerateRandomNumber(len(dtos))]
		if randomPrefix.Pos != "a" && randomPrefix.Pos != "v" {
			continue
		}

		jrp = &GenerateJrpUseCaseOutputDto{
			ID:          0,
			Phrase:      randomPrefix.Lemma + suffix,
			Prefix:      "",
			Suffix:      suffix,
			IsFavorited: 0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		break
	}

	return jrp
}

// RunWithRandom generates a jrp with random prefix and suffix.
func (uc *generateJrpUseCase) RunWithRandom(
	dtos []*GenerateJrpUseCaseInputDto,
) *GenerateJrpUseCaseOutputDto {
	if len(dtos) == 0 {
		return nil
	}

	now := time.Now()
	maxAttempts := len(dtos)

	var jrp *GenerateJrpUseCaseOutputDto = nil
	for i := 0; i < maxAttempts; i++ {
		indexForPrefix := ru.GenerateRandomNumber(len(dtos))
		indexForSuffix := ru.GenerateRandomNumber(len(dtos))

		randomPrefix := dtos[indexForPrefix]
		if randomPrefix.Pos != "a" && randomPrefix.Pos != "v" {
			continue
		}
		randomSuffix := dtos[indexForSuffix]
		if randomSuffix.Pos != "n" {
			continue
		}

		jrp = &GenerateJrpUseCaseOutputDto{
			ID:          0,
			Phrase:      randomPrefix.Lemma + randomSuffix.Lemma,
			Prefix:      "",
			Suffix:      "",
			IsFavorited: 0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		break
	}

	return jrp
}
