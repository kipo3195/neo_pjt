package modules

import (
	appValidation "common/internal/domains/appValidation"
	configuration "common/internal/domains/configuration"
	skin "common/internal/domains/skin"

	"gorm.io/gorm"
)

type Dependencies struct {
	DB                *gorm.DB
	ConfigHashStorage interface{}
}

func InitAppInitModule(dep Dependencies) *handlers.AppInitHandler {
	appValidationUsecase := appValidation.NewUsecase(appValidation.NewAppValidationRepository(db))
	skinUsecase := skin.NewSkinUsecase(skin.NewSkinRepository(db))
	configurationUsecase := configuration.NewUsecase(configuration.NewConfigurationRepository(db))

	svc := services.NewAppInitService(appValidationUsecase, skinUsecase, configurationUsecase)
	return handlers.NewAppInitHander(svc)
}
