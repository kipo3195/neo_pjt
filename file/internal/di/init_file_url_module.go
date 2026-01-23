package di

import (
	"file/internal/application/usecase"
	"file/internal/delivery/handler"
	"file/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type FileUrlModule struct {
	Handler *handler.FileUrlHandler
}

func InitFileUrlModule(db *gorm.DB) *FileUrlModule {

	repo := repository.NewFileUrlRepository(db)
	apiRepo := repository.NewFileUrlApiRepository()
	usecase := usecase.NewFileUrlUsecase(repo, apiRepo)
	handler := handler.NewFileUrlHandler(usecase)

	return &FileUrlModule{
		Handler: handler,
	}
}
