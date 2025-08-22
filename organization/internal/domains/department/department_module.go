package department

import (
	clientHandler "org/internal/domains/department/handlers/client"
	serverHandler "org/internal/domains/department/handlers/server"

	serverRepository "org/internal/domains/department/repositories/client"
	clientRepository "org/internal/domains/department/repositories/server"

	clientUsecase "org/internal/domains/department/usecases/client"
	serverUsecase "org/internal/domains/department/usecases/server"

	"gorm.io/gorm"
)

type DepartmentHandlers struct {
	ClientHandler *clientHandler.DepartmentHandler
	ServerHandler *serverHandler.DepartmentHandler
}

func InitModule(db *gorm.DB) *DepartmentHandlers {

	serverRepository := serverRepository.NewDepartmentRepository(db)
	serverUsecase := serverUsecase.NewDepartmentUsecase(serverRepository)
	serverHandler := serverHandler.NewDepartmentHandler(serverUsecase)

	clientRepository := clientRepository.NewDepartmentRepository(db)
	clientUsecase := clientUsecase.NewDepartmentUsecase(clientRepository)
	clientHandler := clientHandler.NewDepartmentHandler(clientUsecase)

	return &DepartmentHandlers{
		ClientHandler: clientHandler,
		ServerHandler: serverHandler,
	}
}
