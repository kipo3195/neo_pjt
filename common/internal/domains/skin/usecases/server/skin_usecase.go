package server

import (
	repositories "common/internal/domains/skin/repositories/server"
)

type skinUsecase struct {
	repository repositories.SkinRepository
}

type SkinUsecase interface {
}

func NewSkinUsecase(repository repositories.SkinRepository) SkinUsecase {
	return skinUsecase{
		repository: repository,
	}
}
