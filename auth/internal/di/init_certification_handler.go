package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/config"
	"auth/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type CertificationHandler struct {
	Handler *handler.CertificationHandler
}

func InitCertificationHandler(db *gorm.DB, sfg *config.ServerConfig) *CertificationHandler {

	repo := repository.NewCertificationRepository(db)
	usecase := usecase.NewCertificationUsecase(repo)
	handler := handler.NewCertificationHandler(usecase)

	return &CertificationHandler{
		Handler: handler,
	}
}
