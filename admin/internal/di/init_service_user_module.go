package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ServiceUserModule struct {
	Handler *handler.ServiceUserHandler
	Usecase usecase.ServiceUserUsecase
}

func InitServiceUserModule(db *gorm.DB, connector *nats.Conn) ServiceUserModule {

	repo := repository.NewServiceUserRepository(db)
	usecase := usecase.NewServiceUserUsecase(repo, connector)
	handler := handler.NewServiceUserHandler(usecase)

	return ServiceUserModule{
		Handler: handler,
		Usecase: usecase,
	}

}
