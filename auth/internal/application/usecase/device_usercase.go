package usecase

import "auth/internal/domain/device/repository"

type deviceUsecase struct {
	repository repository.DeviceRepository
}

type DeviceUsecase interface {
	CheckDeviceRegistState()
}
