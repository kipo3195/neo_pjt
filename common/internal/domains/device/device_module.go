package device

import (
	serverHandler "common/internal/domains/device/handlers/server"
	"common/internal/domains/device/repositories/serverRepository"
	serverUsecase "common/internal/domains/device/usecases/server"

	"gorm.io/gorm"
)

type DeviceHandlers struct {
	ServerHandler *serverHandler.DeviceHandler
}

func InitModule(db *gorm.DB) *DeviceHandlers {
	serverRepository := serverRepository.NewDeviceRepository(db)
	serverUsecase := serverUsecase.NewDeviceUsecase(serverRepository)
	serverHandler := serverHandler.NewDeviceHandler(serverUsecase)

	return &DeviceHandlers{
		ServerHandler: serverHandler,
	}
}
