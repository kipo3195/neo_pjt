package di

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/repository"
	"batch/internal/infrastructure/storage"
)

type OrgInfoModule struct {
	Usecase usecase.OrgInfoUsecase
}

func InitOrgInfoModule(storage storage.OrgInfoStorage) *OrgInfoModule {

	repo := repository.NewOrgInfoRepository()
	usecase := usecase.NewOrgInfoUsecase(repo, storage)

	return &OrgInfoModule{
		Usecase: usecase,
	}

}
