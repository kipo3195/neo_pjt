package app

import (
	"org/internal/handler"
	"org/internal/infra/repository"
	"org/internal/usecase"

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
