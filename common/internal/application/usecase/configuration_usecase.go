package usecase

import (
	"common/internal/consts"
	"common/internal/domain/configuration/entity"
	"common/internal/domain/configuration/repository"
	"common/internal/infrastructure/storage"
	"log"
)

type configurationUsecase struct {
	repository        repository.ConfigurationRepository
	configHashStorage storage.ConfigHashStorage
}

type ConfigurationUsecase interface {
	CheckConfiguration(configHash string) (bool, error)
	GetWorksConfig() *entity.WorksConfig
}

func NewConfigurationUsecase(repository repository.ConfigurationRepository, configHashStorage storage.ConfigHashStorage) ConfigurationUsecase {
	return &configurationUsecase{
		repository:        repository,
		configHashStorage: configHashStorage,
	}
}

// 변경된 처리
func (r *configurationUsecase) CheckConfiguration(configHash string) (bool, error) {

	skinHash, err := r.configHashStorage.GetConfigHash()
	if err != nil {
		log.Println("서버에 skin hash정보가 없음.")
		return false, consts.ErrSkinHashInvalid
	}

	if skinHash != skinHash {
		log.Println("서버의 skin hash 정보와 다름 server skin hash : ", skinHash)
		return false, consts.ErrSkinHashInvalid
	}

	return true, nil

}

func (r *configurationUsecase) GetWorksConfig() *entity.WorksConfig {

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

	return &entity.WorksConfig{
		TimeZone:   timeZone,
		Language:   language,
		ConfigHash: configHash,
	}
}
