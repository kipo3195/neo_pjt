package org

import (
	clientHandler "org/internal/domains/org/handlers/client"
	serverHandler "org/internal/domains/org/handlers/server"
	"org/internal/infra/storage"

	clientRepository "org/internal/domains/org/repositories/client"
	serverRepository "org/internal/domains/org/repositories/server"

	clientUsecase "org/internal/domains/org/usecases/client"
	serverUsecase "org/internal/domains/org/usecases/server"

	"gorm.io/gorm"
)

type OrgHandlers struct {
	ClientHandler *clientHandler.OrgHandler
	ServerHandler *serverHandler.OrgHandler
}

func InitModule(db *gorm.DB, orgStorage storage.OrgFileStorage) *OrgHandlers {

	serverRepository := serverRepository.NewOrgRepository(db)
	serverUsecase := serverUsecase.NewOrgUsecase(serverRepository, orgStorage)
	serverHandler := serverHandler.NewOrgHandler(serverUsecase)

	clientRepository := clientRepository.NewOrgRepository(db)
	clientUsecase := clientUsecase.NewOrgUsecase(clientRepository, orgStorage)
	clientHandler := clientHandler.NewOrgHandler(clientUsecase)

	return &OrgHandlers{
		ClientHandler: clientHandler,
		ServerHandler: serverHandler,
	}
}
