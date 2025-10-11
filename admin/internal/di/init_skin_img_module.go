package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type SkinImgModule struct {
	Handler *handler.SkinImgHandler
}

func InitSkinImgModule(db *gorm.DB) *SkinImgModule {

	repository := repository.NewSkinImgRepository(db)
	usecase := usecase.NewSkinImgUsecase(repository)
	handler := handler.NewSkinImgHandler(usecase)

	return &SkinImgModule{
		Handler: handler,
	}
}
