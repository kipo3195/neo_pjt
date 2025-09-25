package di

import (
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
	"org/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type DepartmentModule struct {
	Handler *handler.DepartmentHandler
}

func InitDepartmentModule(db *gorm.DB) *DepartmentModule {

	repository := repository.NewDepartmentRepository(db)
	usecase := usecase.NewDepartmentUsecase(repository)
	handler := handler.NewDepartmentHandler(usecase)

	return &DepartmentModule{
		Handler: handler,
	}
}
