package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ServiceUsersModule struct {
	Usecase usecase.ServiceUsersUsecase
}

func InitServiceUsersModule(db *gorm.DB) ServiceUsersModule {
	repository := repository.NewServiceUsersRepository(db)
	usecase := usecase.NewServiceUsersUsecase(repository)

	return ServiceUsersModule{
		Usecase: usecase,
	}
}
