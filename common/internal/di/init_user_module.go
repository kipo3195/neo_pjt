package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"
	"common/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type UserModule struct {
	Usecase usecase.UserUsecase
	Handler *handler.UserHandler
}

func InitUserModule(db *gorm.DB, userStorage storage.UserStorage) *UserModule {

	repo := repository.NewUserRepository(db)
	apiRepository := repository.NewUserAPIRepository()
	usecase := usecase.NewUserUsecase(repo, userStorage, apiRepository)
	handler := handler.NewUserHandler(usecase)

	return &UserModule{
		Usecase: usecase,
		Handler: handler,
	}
}
