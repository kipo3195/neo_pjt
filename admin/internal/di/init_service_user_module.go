package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ServiceUserModule struct {
	Handler *handler.ServiceUserHandler
	Usecase usecase.ServiceUserUsecase
}

func InitServiceUserModule(db *gorm.DB) ServiceUserModule {

	repo := repository.NewServiceUserRepository(db)
	usecase := usecase.NewServiceUserUsecase(repo)
	handler := handler.NewServiceUserHandler(usecase)

	return ServiceUserModule{
		Handler: handler,
		Usecase: usecase,
	}

}
