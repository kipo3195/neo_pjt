package client

import (
	"common/entities"
	"common/internal/consts"
	"common/internal/domains/configuration/dto/client/responseDTO"
	repositories "common/internal/domains/configuration/repositories/client"
	"common/internal/infra/storage"
	"context"
	"fmt"
	"log"
)

type configurationUsecase struct {
	repository        repositories.ConfigurationRepository
	configHashStorage storage.ConfigHashStorage
}

type ConfigurationUsecase interface {
	GetConfigHash(body entities.ConfigHashEntity, ctx context.Context) responseDTO.ConfigHashResult
	CheckConfiguration(configHash string) (bool, error)
}

func NewConfigurationUsecase(repository repositories.ConfigurationRepository, configHashStorage storage.ConfigHashStorage) ConfigurationUsecase {
	return &configurationUsecase{
		repository:        repository,
		configHashStorage: configHashStorage,
	}
}

func (r *configurationUsecase) GetConfigHash(body entities.ConfigHashEntity, ctx context.Context) responseDTO.ConfigHashResult {

	configExist := true
	configSame := false
	skinExist := true
	skinSame := false

	clientConfig := body.ConfigHash
	clientSkin := body.SkinHash
	device := body.Device

	serverConfig, err := r.configHashStorage.GetConfigHash()
	if err != nil {
		log.Println("config 에 대한 hash 정보를 찾을 수 없음.")
		configExist = false
	}

	// 있으면서 동일한
	if configExist && serverConfig == clientConfig {
		configSame = true
	}

	// skin은 skin에서 조회하도록 수정필요
	serverSkin, err := r.skinConigStorage.GetHash(consts.SKIN)
	if err != nil {
		fmt.Printf("%s 에 대한 hash 정보를 찾을 수 없음. \n", device)
		skinExist = false
	}

	// 있으면서 동일한
	if skinExist && serverSkin == clientSkin {
		skinSame = true
	}

	entity := entities.ConfigHashResultEntity{
		ConfigExist: configExist,
		ConfigSame:  configSame,
		SkinExist:   skinExist,
		SkinSame:    skinSame,
	}

	return toConfighashResultDto(entity)
}

func toConfighashResultDto(entity entities.ConfigHashResultEntity) responseDTO.ConfigHashResult {
	return responseDTO.ConfigHashResult{
		ConfigExist: entity.ConfigExist,
		ConfigSame:  entity.ConfigSame,
		SkinExist:   entity.SkinExist,
		SkinSame:    entity.SkinSame,
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
