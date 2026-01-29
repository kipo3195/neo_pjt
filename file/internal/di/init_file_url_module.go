package di

import (
	"file/internal/application/usecase"
	"file/internal/delivery/handler"
	"file/internal/domain/logger"
	"file/internal/infrastructure/cache"
	"file/internal/infrastructure/config"
	"file/internal/infrastructure/repository"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type FileUrlModule struct {
	Handler *handler.FileUrlHandler
}

func InitFileUrlModule(db *gorm.DB, cacheClient *redis.Client, oracleStorageConfig config.OracleStorageConfig, logger logger.Logger) *FileUrlModule {

	cacheStorage := cache.NewFileUrlCache(cacheClient)
	repo := repository.NewFileUrlRepository(db, cacheStorage)
	storageRepo := repository.NewFileUrlStorageRepository(oracleStorageConfig.OciClient, oracleStorageConfig.Namespace, oracleStorageConfig.BucketName)
	usecase := usecase.NewFileUrlUsecase(repo, storageRepo, logger)
	handler := handler.NewFileUrlHandler(usecase)

	return &FileUrlModule{
		Handler: handler,
	}
}
