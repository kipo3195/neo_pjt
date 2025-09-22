package usecase

import "auth/internal/domain/device/repository"

type deviceUsecase struct {
	repo repository.DeviceRepository
}

type DeviceUsecase interface {
}

func NewDeviceUsecase(repo repository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{
		repo: repo,
	}
}
