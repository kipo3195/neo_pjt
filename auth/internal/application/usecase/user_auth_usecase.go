package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/domain/userAuth/entity"
	"auth/internal/domain/userAuth/repository"
	"auth/internal/infrastructure/storage"
	"context"
)

type userAuthUsecase struct {
	repo    repository.UserAuthRepository
	storage storage.UserAuthStorage
}

type UserAuthUsecase interface {
	PutUserAuth(ctx context.Context, input input.UserAuthRegisterInput) string
}

func NewUserAuthUsecase(repo repository.UserAuthRepository, storage storage.UserAuthStorage) UserAuthUsecase {
	return userAuthUsecase{
		repo:    repo,
		storage: storage,
	}
}

func (u userAuthUsecase) PutUserAuth(ctx context.Context, input input.UserAuthRegisterInput) string {

	entity := entity.MakeUserAuthEntity(input.Id, input.Salt, input.AuthHash, input.UserHash)

	err := u.repo.PutUserAuth(ctx, entity)

	if err != nil {
		return "fail"
	}
	return "success"
}
