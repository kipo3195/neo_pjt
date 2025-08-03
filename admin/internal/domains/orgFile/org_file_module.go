package orgFile

import (
	clientHandler "admin/internal/domains/orgFile/handlers/client"
	clientRepository "admin/internal/domains/orgFile/repositories/client"
	clientUsecase "admin/internal/domains/orgFile/usecases/client"

	"gorm.io/gorm"
)

type OrgFileHandlers struct {
	ClientHandler *clientHandler.OrgFileHandler
}

func InitModules(db *gorm.DB) *OrgFileHandlers {

	clientRepository := clientRepository.NewOrgFileRepository(db)
	clientUsecase := clientUsecase.NewOrgFileUsecase(clientRepository)
	clientHandler := clientHandler.NewOrgFileHandler(clientUsecase)

	return &OrgFileHandlers{
		ClientHandler: clientHandler,
	}
}
