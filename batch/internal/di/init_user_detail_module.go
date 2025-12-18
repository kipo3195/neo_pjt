package di

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/repository"
	"batch/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type UserDetailModule struct {
	Usecase usecase.UserDetailUsecase
}

func InitUserDetailModule(db *gorm.DB, storage storage.UserDetailStorage, domain string) *UserDetailModule {

	repo := repository.NewUserDetailRepository(db)
	apiRepo := repository.NewUserDetailApiRepository(domain)
	usecase := usecase.NewUserDetailUsecase(repo, apiRepo)

	return &UserDetailModule{
		Usecase: usecase,
	}

}
