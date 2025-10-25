package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	domainStorage "common/internal/domain/profile/storage"

	"common/internal/infrastructure/repository"
	"common/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ProfileModule struct {
	Usecase usecase.ProfileUsecase
	Handler *handler.ProfileHandler
}

func InitProfileModule(db *gorm.DB, profileStorage domainStorage.ProfileStorage, profileCacheStorage storage.ProfileCacheStorage) *ProfileModule {
	repository := repository.NewProfileRepository(db)
	usecase := usecase.NewProfileUsecase(repository, profileStorage, profileCacheStorage)
	handler := handler.NewProfileHandler(usecase)
	return &ProfileModule{
		Usecase: usecase,
		Handler: handler,
	}
}
