package di

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ExtendDbConnectModule struct {
	Usecase usecase.ExtendDBConnectUsecase
}

func InitExtendDBConnectModule(db *gorm.DB) *ExtendDbConnectModule {

	repo := repository.NewExtendDBConnectRepository(db)
	usecase := usecase.NewExtendDBConnectUsecase(repo)

	return &ExtendDbConnectModule{
		Usecase: usecase,
	}
}
