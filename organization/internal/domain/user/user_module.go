package user

import (
	clientHandler "org/internal/domains/user/handlers/client"
	clientRepository "org/internal/domains/user/repositories/client"
	clientUsecase "org/internal/domains/user/usecases/client"

	"gorm.io/gorm"
)

type UserHandlers struct {
	ClientHandler *clientHandler.UserHandler
}

func InitModule(db *gorm.DB) *UserHandlers {

	clientRepository := clientRepository.NewUserRepository(db)
	clientUsecase := clientUsecase.NewUserUsecase(clientRepository)
	clientHandler := clientHandler.NewUserHandler(clientUsecase)

	return &UserHandlers{
		ClientHandler: clientHandler,
	}
}
