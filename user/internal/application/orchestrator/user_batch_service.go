package orchestrator

import (
	"context"
	"user/internal/application/usecase"
	"user/internal/application/usecase/input"
)

type UserBatchService struct {
	Profile    usecase.ProfileUsecase
	UserDetail usecase.UserDetailUsecase
}

func NewUserBatchService(profile usecase.ProfileUsecase, userDetail usecase.UserDetailUsecase) *UserBatchService {
	return &UserBatchService{
		Profile:    profile,
		UserDetail: userDetail,
	}
}

func (r *UserBatchService) RegisterUserDetailBatch(ctx context.Context, input input.RegistUserDetailBatchInput) error {

	r.UserDetail.RegisterUserDetailBatch(ctx, input)

	return nil
}
