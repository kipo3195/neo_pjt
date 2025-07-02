package usecases

import "admin/repositories"

type commonUsecase struct {
	repo repositories.CommonRepository
}

type CommonUsecase interface {
}

func NewCommonUsecase(repo repositories.CommonRepository) CommonUsecase {
	return &commonUsecase{
		repo: repo,
	}
}
