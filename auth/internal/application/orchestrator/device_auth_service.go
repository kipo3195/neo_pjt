package orchestrator

import "auth/internal/application/usecase"

type DeviceAuthService struct {
	Token  usecase.TokenUsecase
	Device usecase.DeviceUsecase
}

func NewDeviceAuthService(t usecase.TokenUsecase, d usecase.DeviceUsecase) *DeviceAuthService {
	return &DeviceAuthService{
		Token:  t,
		Device: d,
	}

}
