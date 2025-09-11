package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type WorksInfoModule struct {
	Handler *handler.WorksInfoHandler
	Usecase usecase.WorksInfoUsecase
}

func InitWorksInfoHandler(db *gorm.DB) *WorksInfoModule {
	repository := repository.NewWorksInfoRepository(db)
	usecase := usecase.NewWorksInfoUsecase(repository)
	handler := handler.NewWorksInfoHandler(usecase)

	return &WorksInfoModule{
		Handler: handler,
		Usecase: usecase,
	}
}
