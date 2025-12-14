package di

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/repository"
	"batch/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type OrgInfoModule struct {
	Usecase usecase.OrgInfoUsecase
}

func InitOrgInfoModule(db *gorm.DB, storage storage.OrgInfoStorage, domain string) *OrgInfoModule {

	repo := repository.NewOrgInfoRepository(db)
	apiRepo := repository.NewOrgInfoApiRepository(domain)
	usecase := usecase.NewOrgInfoUsecase(repo, apiRepo, storage)

	return &OrgInfoModule{
		Usecase: usecase,
	}

}
