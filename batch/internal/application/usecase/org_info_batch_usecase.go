package usecase

import (
	"batch/internal/domain/orgInfo/repository"
	"batch/internal/infrastructure/storage"
)

type orgInfoBatchUsecase struct {
	orgInfoBatchStorage storage.OrgInfoBatchStorage
	repo                repository.OrgInfoBatchRepository
}

type OrgInfoBatchUsecase interface {
}

func NewOrgInfoBatchUsecase(repo repository.OrgInfoBatchRepository, orgInfoBatchStorage storage.OrgInfoBatchStorage) OrgInfoBatchUsecase {

	return &orgInfoBatchUsecase{
		orgInfoBatchStorage: orgInfoBatchStorage,
		repo:                repo,
	}
}
