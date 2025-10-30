package di

import (
	"user/internal/application/usecase"
	"user/internal/delivery/handler"
	"user/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type UserDetailModule struct {
	Usecase usecase.UserDetailUsecase
	Handler *handler.UserDetailHandler
}

func InitUserDetailModule(db *gorm.DB) *UserDetailModule {
	repository := repository.NewUserDetailRepository(db)
	usecase := usecase.NewUserDatailUsecase(repository)
	handler := handler.NewUserDetailHandler(usecase)

	return &UserDetailModule{
		Usecase: usecase,
		Handler: handler,
	}

}
