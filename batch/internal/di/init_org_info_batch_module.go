package di

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/repository"
	"batch/internal/infrastructure/storage"
)

type OrgInfoBatchModule struct {
	Usecase usecase.OrgInfoBatchUsecase
}

func InitOrgInfoBatchModule(storage storage.OrgInfoBatchStorage) *OrgInfoBatchModule {

	repo := repository.NewOrgInfoBatchRepository()
	usecase := usecase.NewOrgInfoBatchUsecase(repo, storage)

	return &OrgInfoBatchModule{
		Usecase: usecase,
	}

}
