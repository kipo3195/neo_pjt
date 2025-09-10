package device

import (
	serverHandler "common/internal/domains/worksInfo/handlers/server"
	"common/internal/domains/worksInfo/repositories/serverRepository"
	serverUsecase "common/internal/domains/worksInfo/usecases/server"

	"gorm.io/gorm"
)

type WorksInfoHandler struct {
	ServerHandler *serverHandler.WorksInfoHandler
}

func InitModule(db *gorm.DB) *WorksInfoHandler {
	serverRepository := serverRepository.NewWorksInfoRepository(db)
	serverUsecase := serverUsecase.NewWorksInfoUsecase(serverRepository)
	serverHandler := serverHandler.NewWorksInfoHandler(serverUsecase)

	return &WorksInfoHandler{
		ServerHandler: serverHandler,
	}
}
