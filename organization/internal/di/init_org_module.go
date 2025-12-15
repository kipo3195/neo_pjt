package di

import (
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
	"org/internal/infrastructure/repository"
	"org/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type OrgModule struct {
	Handler *handler.OrgHandler
	Usecase usecase.OrgUsecase
}

func InitOrgModule(db *gorm.DB, orgFileStorage storage.OrgFileStorage, orgStorage storage.OrgStorage) *OrgModule {

	repository := repository.NewOrgRepository(db)
	usecase := usecase.NewOrgUsecase(repository, orgFileStorage, orgStorage)
	handler := handler.NewOrgHandler(usecase)

	return &OrgModule{
		Handler: handler,
		Usecase: usecase,
	}
}
