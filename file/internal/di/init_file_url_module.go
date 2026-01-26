package di

import (
	"file/internal/application/usecase"
	"file/internal/delivery/handler"
	"file/internal/domain/logger"
	"file/internal/infrastructure/config"
	"file/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type FileUrlModule struct {
	Handler *handler.FileUrlHandler
}

func InitFileUrlModule(db *gorm.DB, oracleStorageConfig config.OracleStorageConfig, logger logger.Logger) *FileUrlModule {

	repo := repository.NewFileUrlRepository(db)
	storageRepo := repository.NewFileUrlStorageRepository(oracleStorageConfig.OciClient, oracleStorageConfig.Namespace, oracleStorageConfig.BucketName)
	usecase := usecase.NewFileUrlUsecase(repo, storageRepo, logger)
	handler := handler.NewFileUrlHandler(usecase)

	return &FileUrlModule{
		Handler: handler,
	}
}
