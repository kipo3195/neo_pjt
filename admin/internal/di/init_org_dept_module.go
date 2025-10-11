package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type OrgDeptModule struct {
	Handler *handler.OrgDeptHandler
}

func InitOrgDeptModule(db *gorm.DB) *OrgDeptModule {

	repository := repository.NewOrgDeptRepository(db)
	usecase := usecase.NewOrgDeptsUsecase(repository)
	handler := handler.NewOrgDeptsHandler(usecase)

	return &OrgDeptModule{
		Handler: handler,
	}
}
