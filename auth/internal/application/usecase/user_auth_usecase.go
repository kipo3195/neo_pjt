package usecase

import "auth/internal/domain/userAuth/repository"

type userAuthUsecase struct {
	repo repository.UserAuthRepository
}

type UserAuthUsecase interface {
}

func NewUserAuthUsecase(repo repository.UserAuthRepository) UserAuthUsecase {
	return userAuthUsecase{
		repo: repo,
	}
}
