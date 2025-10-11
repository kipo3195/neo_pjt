package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type OrgDeptUserModule struct {
	Handler *handler.OrgDeptUserHandler
}

func InitOrgDeptUserModule(db *gorm.DB) *OrgDeptUserModule {

	repository := repository.NewOrgDeptUserRepository(db)
	usecase := usecase.NewOrgDeptUsersUsecase(repository)
	handler := handler.NewOrgDeptUsersHandler(usecase)

	return &OrgDeptUserModule{
		Handler: handler,
	}
}
