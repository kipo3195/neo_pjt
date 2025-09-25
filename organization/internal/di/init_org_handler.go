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
}

func InitOrgModule(db *gorm.DB, orgStorage storage.OrgFileStorage) *OrgModule {

	repository := repository.NewOrgRepository(db)
	usecase := usecase.NewOrgUsecase(repository, orgStorage)
	handler := handler.NewOrgHandler(usecase)

	return &OrgModule{
		Handler: handler,
	}
}
