package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type AppTokenModule struct {
	Handler *handler.AppTokenHandler
	Usecase usecase.AppTokenUsecase
}

func InitAppTokenModule(db *gorm.DB) *AppTokenModule {

	repository := repository.NewAppTokenRepository(db)
	usecase := usecase.NewAppTokenUsecase(repository)
	handler := handler.NewAppTokenHandler(usecase)

	return &AppTokenModule{
		Handler: handler,
		Usecase: usecase,
	}
}
