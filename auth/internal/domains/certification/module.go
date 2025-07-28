package certification

import (
	clHandler "auth/internal/domains/certification/handlers/client"
	clRepo "auth/internal/domains/certification/repositories/client"
	clUsecase "auth/internal/domains/certification/usecases/client"
	"auth/pkg/config"

	"gorm.io/gorm"
)

func InitCertificationModule(db *gorm.DB, jwtCfg *config.JWTConfig) *clHandler.CertificationHandler {
	repository := clRepo.NewCertificationRepository(db)
	usecase := clUsecase.NewCertificationUsecase(repository, jwtCfg)
	handler := clHandler.NewCertificationHandler(usecase)

	return handler
}
