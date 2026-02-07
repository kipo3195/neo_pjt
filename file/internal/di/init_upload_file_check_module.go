package di

import (
	"file/internal/application/usecase"
	"file/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type UploadFileCheckModule struct {
	Usecase usecase.UploadFileCheckUsecase
}

func InitUploadFileCheckModule(db *gorm.DB) UploadFileCheckModule {

	repository := repository.NewUploadFileCheckRepository(db)
	usecase := usecase.NewUploadFileCheckUsecase(repository)

	return UploadFileCheckModule{
		Usecase: usecase,
	}
}
