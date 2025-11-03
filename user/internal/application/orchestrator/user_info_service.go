package orchestrator

import "user/internal/application/usecase"

type UserInfoService struct {
	Profile    usecase.ProfileUsecase
	UserDetail usecase.UserDetailUsecase
}

func NewUserInfoService(p usecase.ProfileUsecase, u usecase.UserDetailUsecase) *UserInfoService {
	return &UserInfoService{
		Profile:    p,
		UserDetail: u,
	}
}
