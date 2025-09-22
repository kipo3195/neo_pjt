package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type DeviceModule struct {
	Handler *handler.DeviceHandler
	Usecase usecase.DeviceUsecase
}

func InitDeviceModule(db *gorm.DB) DeviceModule {

	repo := repository.NewDeviceRepository(db)
	usecase := usecase.NewDeviceUsecase(repo)
	handler := handler.NewDeviceHandler(usecase)
	return DeviceModule{
		Handler: handler,
		Usecase: usecase,
	}
}
