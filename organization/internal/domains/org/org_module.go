package org

import (
	clientHandler "org/internal/domains/org/handlers/client"
	serverHandler "org/internal/domains/org/handlers/server"

	serverRepository "org/internal/domains/org/repositories/client"
	clientRepository "org/internal/domains/org/repositories/server"

	clientUsecase "org/internal/domains/org/usecases/client"
	serverUsecase "org/internal/domains/org/usecases/server"

	"gorm.io/gorm"
)

type OrgHandlers struct {
	ClientHandler *clientHandler.OrgHandler
	ServerHandler *serverHandler.OrgHandler
}

func InitModule(db *gorm.DB) *OrgHandlers {

	serverRepository := serverRepository.NewOrgRepository(db)
	serverUsecase := serverUsecase.NewOrgUsecase(serverRepository)
	serverHandler := serverHandler.NewOrgHandler(serverUsecase)

	clientRepository := clientRepository.NewOrgRepository(db)
	clientUsecase := clientUsecase.NewOrgUsecase(clientRepository)
	clientHandler := clientHandler.NewOrgHandler(clientUsecase)

	return &OrgHandlers{
		ClientHandler: clientHandler,
		ServerHandler: serverHandler,
	}
}
