package orchestrator

import "auth/internal/application/usecase"

type UserAuthService struct {
	UserAuth usecase.UserAuthUsecase
	Device   usecase.DeviceUsecase
	Token    usecase.TokenUsecase
}

func NewUserAuthService(u usecase.UserAuthUsecase, d usecase.DeviceUsecase, t usecase.TokenUsecase) *UserAuthService {
	return &UserAuthService{
		UserAuth: u,
		Device:   d,
		Token:    t,
	}
}
