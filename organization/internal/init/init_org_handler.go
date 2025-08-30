package init

import (
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
	"org/internal/infrastructure/repository"
	"org/internal/infrastructure/storage"

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
