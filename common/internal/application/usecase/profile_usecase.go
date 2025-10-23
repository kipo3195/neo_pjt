package usecase

import (
	"common/internal/domain/profile/repository"
	"common/internal/domain/profile/storage"
)

type profileUsecase struct {
	repository     repository.ProfileRepository
	profileStorage storage.ProfileStorage
}

type ProfileUsecase interface {
}

func NewProfileUsecase(repository repository.ProfileRepository, profileStorage storage.ProfileStorage) ProfileUsecase {
	return profileUsecase{
		repository:     repository,
		profileStorage: profileStorage,
	}
}
