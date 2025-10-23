package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/domain/profile/storage"
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ProfileModule struct {
	Usecase usecase.ProfileUsecase
	Handler *handler.ProfileHandler
}

func InitProfileModule(db *gorm.DB, profileStorage storage.ProfileStorage) *ProfileModule {
	repository := repository.NewProfileRepository(db)
	usecase := usecase.NewProfileUsecase(repository, profileStorage)
	handler := handler.NewProfileHandler(usecase)
	return &ProfileModule{
		Usecase: usecase,
		Handler: handler,
	}
}
