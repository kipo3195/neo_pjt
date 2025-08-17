package server

import (
	"common/internal/consts"
	"common/internal/domains/worksInfo/dto/server/requestDTO"
	"common/internal/domains/worksInfo/entities"
	"common/internal/domains/worksInfo/repositories/serverRepository"
)

type worksInfoUsecase struct {
	repository serverRepository.WorksInfoRepository
}

type WorksInfoUsecase interface {
	GetConnectInfo(body *requestDTO.ConnectInfoRequest) (*entities.ConnectInfo, error)
}

func NewWorksInfoUsecase(repository serverRepository.WorksInfoRepository) WorksInfoUsecase {
	return &worksInfoUsecase{
		repository: repository,
	}
}

func (u *worksInfoUsecase) GetConnectInfo(body *requestDTO.ConnectInfoRequest) (*entities.ConnectInfo, error) {

	connectInfo, err := u.repository.GetConnectInfo(body.WorksCode)

	if err != nil {
		return nil, consts.ErrDB
	}

	return connectInfo, nil

}
