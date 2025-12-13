package usecase

import "batch/internal/domain/extendDbConnect/repository"

type extendDBConnectUsecase struct {
}

type ExtendDBConnectUsecase interface {
}

func NewExtendDBConnectUsecase(repo repository.ExtendDBConnectRepository) ExtendDBConnectUsecase {
	return &extendDBConnectUsecase{}
}
