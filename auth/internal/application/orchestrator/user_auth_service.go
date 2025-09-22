package orchestrator

import "auth/internal/application/usecase"

type UserAuthService struct {
	UserAuth usecase.UserAuthUsecase
	Device   usecase.DeviceUsecase
}

func NewUserAuthService(u usecase.UserAuthUsecase, d usecase.DeviceUsecase) *UserAuthService {
	return &UserAuthService{
		UserAuth: u,
		Device:   d,
	}
}
