package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type OrgFileModule struct {
	Handler *handler.OrgFileHandler
}

func InitOrgFileModule(db *gorm.DB) *OrgFileModule {

	repository := repository.NewOrgFileRepository(db)
	usecase := usecase.NewOrgFileUsecase(repository)
	handler := handler.NewOrgFileHandler(usecase)

	return &OrgFileModule{
		Handler: handler,
	}
}
