package usecase

import (
	"batch/internal/domain/userDetail/repository"
	"context"
	"log"
)

type userDetailUsecase struct {
	repo    repository.UserDetailRepository
	apiRepo repository.UserDetailApiRepository
}

type UserDetailUsecase interface {
	SendUserDetailToUser(ctx context.Context, org string) error
}

func NewUserDetailUsecase(repo repository.UserDetailRepository, apiRepo repository.UserDetailApiRepository) UserDetailUsecase {
	return &userDetailUsecase{
		repo:    repo,
		apiRepo: apiRepo,
	}
}

func (r *userDetailUsecase) SendUserDetailToUser(ctx context.Context, org string) error {

	// 현재 데이터 조회
	userDetail, err := r.repo.GetUserDetail(ctx, org)

	if err != nil {
		return err
	}

	log.Println("SendUserDetailToUser 현재 데이터 :", userDetail)

	return nil
}
