package usecases

import (
	"admin/repositories"
)

type authUsecase struct {
	repo repositories.AuthRepository
}

type AuthUsecase interface {
}

func NewAuthUsecase(repo repositories.AuthRepository) AuthUsecase {
	return &authUsecase{repo: repo}
}
