package orgDepts

import (
	clientHandler "admin/internal/domains/orgDepts/handlers/client"
	clientRepository "admin/internal/domains/orgDepts/repositories/client"
	clientUsecase "admin/internal/domains/orgDepts/usecases/client"

	"gorm.io/gorm"
)

type OrgDeptsHandlers struct {
	ClientHandler *clientHandler.OrgDeptsHandler
}

func InitModules(db *gorm.DB) *OrgDeptsHandlers {

	clientRepository := clientRepository.NewOrgDeptsRepository(db)
	clientUsecase := clientUsecase.NewOrgDeptsUsecase(clientRepository)
	clientHandler := clientHandler.NewOrgDeptsHandler(clientUsecase)

	return &OrgDeptsHandlers{
		ClientHandler: clientHandler,
	}
}
