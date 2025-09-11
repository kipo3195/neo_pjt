package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"
	"common/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type SkinModule struct {
	Usecase usecase.SkinUsecase
	Handler *handler.SkinHandler
}

func InitSkinModule(db *gorm.DB, skinStorage storage.SkinStorage) *SkinModule {
	repository := repository.NewSkinRepository(db)
	usecase := usecase.NewSkinUsecase(repository, skinStorage) // storage 필요
	handler := handler.NewSkinHandler(usecase)

	return &SkinModule{
		Usecase: usecase,
		Handler: handler,
	}
}
