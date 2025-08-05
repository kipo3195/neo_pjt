package device

import "gorm.io/gorm"

type DeviceHandlers struct {
	ServerHandler *serverHandler.DeviceHandlers
}

func InitModule(db *gorm.DB) *DeviceHandlers {
	serverRepository := serverRepository.NewDeviceRepository(db)
	serverUsecase := serverUsecase.NewDeviceUsecase(serverRepository)
	serverHandler := serverHandler.NewDeviceHandler(serverUsecase)

	return &DeviceHandlers{
		ServerHandler: serverHandler,
	}
}
