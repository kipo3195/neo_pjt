package di

import (
	"user/internal/application/usecase"
	"user/internal/delivery/handler"
	domainStorage "user/internal/domain/profile/storage"

	"user/internal/infrastructure/repository"
	"user/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ProfileModule struct {
	Usecase usecase.ProfileUsecase
	Handler *handler.ProfileHandler
}

func InitProfileModule(db *gorm.DB, profileStorage domainStorage.ProfileStorage, profileCacheStorage storage.ProfileCacheStorage, userInfoServiceStorage storage.UserInfoServiceStorage) *ProfileModule {
	repository := repository.NewProfileRepository(db)
	usecase := usecase.NewProfileUsecase(repository, profileStorage, profileCacheStorage, userInfoServiceStorage)
	handler := handler.NewProfileHandler(usecase)
	return &ProfileModule{
		Usecase: usecase,
		Handler: handler,
	}
}
