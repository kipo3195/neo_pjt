package orchestrator

import "auth/internal/application/usecase"

type UserAuthTokenService struct {
	UserAuth usecase.UserAuthUsecase
	Token    usecase.TokenUsecase
}

func NewUserAuthTokenService(u usecase.UserAuthUsecase, t usecase.TokenUsecase) *UserAuthTokenService {
	return &UserAuthTokenService{
		UserAuth: u,
		Token:    t,
	}
}
