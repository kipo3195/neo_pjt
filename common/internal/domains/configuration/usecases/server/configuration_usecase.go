package server

import (
	"common/internal/domains/configuration/entities"
	"common/internal/infra/storage"
	"log"
)

type configurationUsecase struct {
	configHashStorage storage.ConfigHashStorage
}

type ConfigurationUsecase interface {
	GetWorksConfig() *entities.WorksConfig
}

func NewConfigurationUsecase(configHashStorage storage.ConfigHashStorage) ConfigurationUsecase {
	return &configurationUsecase{
		configHashStorage: configHashStorage,
	}
}

func (r *configurationUsecase) GetWorksConfig() *entities.WorksConfig {

	configHash, err := r.configHashStorage.GetConfigHash()
	if err != nil {
		log.Println("[GetWorksInfo] configHash is invalid")
	}

	timeZone, err := r.configHashStorage.GetWorkConfig("timeZone")
	if err != nil {
		log.Println("[GetWorksInfo] timeZone is invalid")
	}

	language, err := r.configHashStorage.GetWorkConfig("language")
	if err != nil {
		log.Println("[GetWorksInfo] language is invalid")
	}

	return &entities.WorksConfig{
		TimeZone:   timeZone,
		Language:   language,
		ConfigHash: configHash,
	}
}
