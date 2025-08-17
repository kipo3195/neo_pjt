package server

import (
	"common/internal/consts"
	"common/internal/domains/device/dto/server/requestDTO"
	"common/internal/domains/device/entities"
	"common/internal/domains/device/repositories/serverRepository"
)

type deviceUsecase struct {
	repository serverRepository.DeviceRepository
}

type DeviceUsecase interface {
	GetConnectInfo(body *requestDTO.DeviceInitRequest) (*entities.ConnectInfo, error)
}

func NewDeviceUsecase(repository serverRepository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{
		repository: repository,
	}
}

func (u *deviceUsecase) GetConnectInfo(body *requestDTO.DeviceInitRequest) (*entities.ConnectInfo, error) {

	connectInfo, err := u.repository.GetConnectInfo(body.WorksCode)

	if err != nil {
		return nil, consts.ErrDB
	}

	return connectInfo, nil

}
