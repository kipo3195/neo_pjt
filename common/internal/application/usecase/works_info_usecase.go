package usecase

import (
	"common/internal/consts"
	"common/internal/delivery/dto/worksInfo"
	"common/internal/domain/worksInfo/entity"
	"common/internal/domain/worksInfo/repository"
)

type worksInfoUsecase struct {
	repository repository.WorksInfoRepository
}

type WorksInfoUsecase interface {
	GetConnectInfo(body *worksInfo.ConnectInfoRequest) (*entity.ConnectInfo, error)
}

func NewWorksInfoUsecase(repository repository.WorksInfoRepository) WorksInfoUsecase {
	return &worksInfoUsecase{
		repository: repository,
	}
}

func (u *worksInfoUsecase) GetConnectInfo(body *worksInfo.ConnectInfoRequest) (*entity.ConnectInfo, error) {

	connectInfo, err := u.repository.GetConnectInfo(body.WorksCode)

	if err != nil {
		return nil, consts.ErrDB
	}

	return connectInfo, nil

}
