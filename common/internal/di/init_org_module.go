package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"
	"common/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type OrgModule struct {
	Handler *handler.OrgHandler
	Usecase usecase.OrgUsecase
}

func InitOrgModule(db *gorm.DB, orgStorage storage.OrgStorage) OrgModule {

	repo := repository.NewOrgRepository(db)
	usecase := usecase.NewOrgUsecase(repo, orgStorage)
	handler := handler.NewOrgHandler(usecase)

	return OrgModule{
		Handler: handler,
		Usecase: usecase,
	}
}
