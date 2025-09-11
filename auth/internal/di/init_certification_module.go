package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/config"
	"auth/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type CertificationModule struct {
	Handler *handler.CertificationHandler
}

func InitCertificationModule(db *gorm.DB, sfg *config.ServerConfig) *CertificationModule {

	repo := repository.NewCertificationRepository(db)
	usecase := usecase.NewCertificationUsecase(repo)
	handler := handler.NewCertificationHandler(usecase)

	return &CertificationModule{
		Handler: handler,
	}
}
