package orgDeptUsers

import (
	clientHandler "admin/internal/domains/orgDeptUsers/handlers/client"
	clientRepository "admin/internal/domains/orgDeptUsers/repositories/client"
	clientUsecase "admin/internal/domains/orgDeptUsers/usecases/client"

	"gorm.io/gorm"
)

type OrgDeptUsersHandlers struct {
	ClientHandler *clientHandler.OrgDeptUsersHandler
}

func InitModules(db *gorm.DB) *OrgDeptUsersHandlers {

	clientRepository := clientRepository.NewOrgDeptUsersRepository(db)
	clientUsecase := clientUsecase.NewOrgDeptUsersUsecase(clientRepository)
	clientHandler := clientHandler.NewOrgDeptUsersHandler(clientUsecase)

	return &OrgDeptUsersHandlers{
		ClientHandler: clientHandler,
	}
}
