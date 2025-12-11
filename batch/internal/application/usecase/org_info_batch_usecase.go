package usecase

import (
	"batch/internal/domain/orgInfo/repository"
	"batch/internal/infrastructure/storage"
	"log"
	"time"
)

type orgInfoBatchUsecase struct {
	orgInfoBatchStorage storage.OrgInfoBatchStorage
	repo                repository.OrgInfoBatchRepository
}

type OrgInfoBatchUsecase interface {
	Run() error
}

func NewOrgInfoBatchUsecase(repo repository.OrgInfoBatchRepository, orgInfoBatchStorage storage.OrgInfoBatchStorage) OrgInfoBatchUsecase {

	return &orgInfoBatchUsecase{
		orgInfoBatchStorage: orgInfoBatchStorage,
		repo:                repo,
	}
}

func (r *orgInfoBatchUsecase) Run() error {

	log.Println("Current Time:", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}
