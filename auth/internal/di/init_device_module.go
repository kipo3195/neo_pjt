package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type DeviceModule struct {
	Handler *handler.DeviceHandler
	Usecase usecase.DeviceUsecase
}

func InitDeviceModule(db *gorm.DB, deviceStorage storage.DeviceStorage, authTokenStorage storage.AuthTokenStorage, accessHash string, refreshHash string) DeviceModule {

	repo := repository.NewDeviceRepository(db)
	usecase := usecase.NewDeviceUsecase(repo, deviceStorage, authTokenStorage, accessHash, refreshHash)
	handler := handler.NewDeviceHandler(usecase)
	return DeviceModule{
		Handler: handler,
		Usecase: usecase,
	}
}
