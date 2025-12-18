package usecase

import "batch/internal/domain/extendDbConnect/repository"

type extendDBConnectUsecase struct {
}

// 설정을 어디서 주입해야되는지 점검..
type ExtendDBConnectUsecase interface {
	GetOrgInfo() error
	GetUserDetail() error
}

func NewExtendDBConnectUsecase(repo repository.ExtendDBConnectRepository) ExtendDBConnectUsecase {
	return &extendDBConnectUsecase{}
}

func (r *extendDBConnectUsecase) GetOrgInfo() error {
	return nil
}

func (r *extendDBConnectUsecase) GetUserDetail() error {
	return nil
}
