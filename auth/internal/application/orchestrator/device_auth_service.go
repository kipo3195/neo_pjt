package orchestrator

import "auth/internal/application/usecase"

type DeviceAuthService struct {
	Token  usecase.TokenUsecase
	Device usecase.DeviceUsecase
	Otp    usecase.OtpUsecase
}

func NewDeviceAuthService(t usecase.TokenUsecase, d usecase.DeviceUsecase, o usecase.OtpUsecase) *DeviceAuthService {
	return &DeviceAuthService{
		Token:  t,
		Device: d,
		Otp:    o,
	}

}
