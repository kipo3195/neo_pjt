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

func InitOrgInfoModule(db *gorm.DB, storage storage.OrgInfoStorage) *OrgInfoModule {

	repo := repository.NewOrgInfoRepository(db)
	usecase := usecase.NewOrgInfoUsecase(repo, storage)

	return &OrgInfoModule{
		Usecase: usecase,
	}

}
