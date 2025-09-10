package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type SkinHandler struct {
	Handler *handler.SkinHandler
}

func InitSkinHandler(db *gorm.DB) *SkinHandler {
	repository := repository.NewSkinRepository(db)
	usecase := usecase.NewSkinUsecase(repository, nil) // storage 필요
	handler := handler.NewSkinHandler(usecase)

	return &SkinHandler{
		Handler: handler,
	}
}
