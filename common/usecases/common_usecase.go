package usecases

import (
	clDto "common/dto/client"
	"common/entities"
	"common/infra/storage"
	"common/repositories"
	"context"
	"fmt"
)

type commonUsecase struct {
	repo              repositories.CommonRepository
	configHashStorage storage.ConfigHashStorage
}

type CommonUsecase interface {
	InitConfigHash() error
	GetConfigHash(body entities.ConfigHashEntity, ctx context.Context) clDto.ConfigHashResult
}

func NewCommonUsecase(repo repositories.CommonRepository, configHashStorage storage.ConfigHashStorage) CommonUsecase {
	return &commonUsecase{
		repo:              repo,
		configHashStorage: configHashStorage,
	}
}

func toWorksConfigEntity(worksCode string, device string) entities.GetWorksConfig {

	return entities.GetWorksConfig{
		WorksCode: worksCode,
		Device:    device,
	}

}

func (r *commonUsecase) InitConfigHash() error {

	// skin 정보 init
	skinHashs, err := r.repo.GetSkinHashs()
	if err != nil {
		return err
	}

	// 없어도 에러는 아님
	for _, m := range skinHashs {
		r.configHashStorage.SaveConfigHash(m.Device, m.SkinHash)
	}

	// config 정보 init
	config, err := r.repo.GetConfig()
	if err != nil {
		return err
	}

	// 없어도 에러는 아님
	if config.ConfigHash != "" && config.Device != "" {
		r.configHashStorage.SaveConfigHash(config.Device, config.ConfigHash)
	}

	return nil
}

func (r *commonUsecase) GetConfigHash(body entities.ConfigHashEntity, ctx context.Context) clDto.ConfigHashResult {

	configExist := true
	configSame := false
	skinExist := true
	skinSame := false

	clientConfig := body.ConfigHash
	clientSkin := body.SkinHash
	device := body.Device

	serverConfig, err := r.configHashStorage.GetConfigHash("config")
	if err != nil {
		fmt.Println("config 에 대한 hash 정보를 찾을 수 없음.")
		configExist = false
	}

	// 있으면서 동일한
	if configExist && serverConfig == clientConfig {
		configSame = true
	}

	serverSkin, err := r.configHashStorage.GetConfigHash(device)
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

func toConfighashResultDto(entity entities.ConfigHashResultEntity) clDto.ConfigHashResult {
	return clDto.ConfigHashResult{
		ConfigExist: entity.ConfigExist,
		ConfigSame:  entity.ConfigSame,
		SkinExist:   entity.SkinExist,
		SkinSame:    entity.SkinSame,
	}
}
