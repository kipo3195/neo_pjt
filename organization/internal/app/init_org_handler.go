package app

import (
	"org/internal/handler"
	"org/internal/infra/repository"
	"org/internal/infra/storage"
	"org/internal/usecase"

	"gorm.io/gorm"
)

type OrgHandler struct {
	Handler *handler.OrgHandler
}

func InitOrgHandler(db *gorm.DB, orgStorage storage.OrgFileStorage) *OrgHandler {

	repository := repository.NewOrgRepository(db)
	usecase := usecase.NewOrgUsecase(repository, orgStorage)
	handler := handler.NewOrgHandler(usecase)

	return &OrgHandler{
		Handler: handler,
	}
}
