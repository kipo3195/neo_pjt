package usecase

import (
	"common/internal/application/usecase/input"
	"common/internal/consts"
	"common/internal/domain/worksInfo/entity"
	"common/internal/domain/worksInfo/repository"
)

type worksInfoUsecase struct {
	repository repository.WorksInfoRepository
}

type WorksInfoUsecase interface {
	GetConnectInfo(input *input.ConnectInfoInput) (*entity.ConnectInfo, error)
}

func NewWorksInfoUsecase(repository repository.WorksInfoRepository) WorksInfoUsecase {
	return &worksInfoUsecase{
		repository: repository,
	}
}

func (u *worksInfoUsecase) GetConnectInfo(input *input.ConnectInfoInput) (*entity.ConnectInfo, error) {

	connectInfo, err := u.repository.GetConnectInfo(input.WorksCode)

	if err != nil {
		return nil, consts.ErrDB
	}

	return connectInfo, nil

}
