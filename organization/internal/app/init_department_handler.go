package app

import (
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
	"org/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type DepartmentHandler struct {
	Handler *handler.DepartmentHandler
}

func InitDepartmentHandler(db *gorm.DB) *DepartmentHandler {

	repository := repository.NewDepartmentRepository(db)
	usecase := usecase.NewDepartmentUsecase(repository)
	handler := handler.NewDepartmentHandler(usecase)

	return &DepartmentHandler{
		Handler: handler,
	}
}
