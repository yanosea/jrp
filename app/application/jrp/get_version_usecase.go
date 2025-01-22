package jrp

import ()

// getVersionUseCase is a struct that contains the use case of getting the version.
type getVersionUseCase struct{}

// NewGetVersionUseCase returns a new instance of the GetVersionUseCase struct.
func NewGetVersionUseCase() *getVersionUseCase {
	return &getVersionUseCase{}
}

// GetVersionUseCaseOutputDto is a DTO struct that contains the output data of the GetVersionUseCase.
type GetVersionUseCaseOutputDto struct {
	Version string
}

// Run returns the output of the GetVersionUseCase.
func (uc *getVersionUseCase) Run(version string) *GetVersionUseCaseOutputDto {
	return &GetVersionUseCaseOutputDto{
		Version: version,
	}
}
