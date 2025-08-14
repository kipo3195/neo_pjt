package client

import (
	"common/internal/consts"
	repositories "common/internal/domains/configuration/repositories/client"
	"common/internal/infra/storage"
	"log"
)

type configurationUsecase struct {
	repository        repositories.ConfigurationRepository
	configHashStorage storage.ConfigHashStorage
}

type ConfigurationUsecase interface {
	CheckConfiguration(configHash string) (bool, error)
}

func NewConfigurationUsecase(repository repositories.ConfigurationRepository, configHashStorage storage.ConfigHashStorage) ConfigurationUsecase {
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
