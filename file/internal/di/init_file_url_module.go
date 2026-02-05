package di

import (
	"file/internal/adapter/http/handler"
	"file/internal/application/usecase"
	"file/internal/domain/logger"
	"file/internal/infrastructure/config"
	"file/internal/infrastructure/persistence/cacheStorage"
	"file/internal/infrastructure/persistence/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type FileUrlModule struct {
	Handler *handler.FileUrlHandler
}

func InitFileUrlModule(db *gorm.DB, cacheClient *redis.ClusterClient, oracleStorageConfig config.OracleStorageConfig, logger logger.Logger) *FileUrlModule {

	cacheStorage := cacheStorage.NewFileUrlCache(cacheClient)
	repo := repository.NewFileUrlRepository(db, cacheStorage)
	storageRepo := repository.NewFileUrlStorageRepository(oracleStorageConfig.OciClient, oracleStorageConfig.Namespace, oracleStorageConfig.BucketName)
	usecase := usecase.NewFileUrlUsecase(repo, storageRepo, logger)
	handler := handler.NewFileUrlHandler(usecase)

	return &FileUrlModule{
		Handler: handler,
	}
}
