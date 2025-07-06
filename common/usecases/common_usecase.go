package usecases

import (
	"common/consts"
	clDto "common/dto/client"
	"common/entities"
	"common/infra/storage"
	"common/repositories"
	"context"
	"fmt"
	"os"
)

type commonUsecase struct {
	repo          repositories.CommonRepository
	configStorage storage.ConfigStorage
}

type CommonUsecase interface {
	InitConfigHash() error
	GetConfigHash(body entities.ConfigHashEntity, ctx context.Context) clDto.ConfigHashResult
	GetSkinImg(ctx context.Context, dto clDto.GetSkinImgRequest) (*os.File, error)
}

func NewCommonUsecase(repo repositories.CommonRepository, configHashStorage storage.ConfigStorage) CommonUsecase {
	return &commonUsecase{
		repo:          repo,
		configStorage: configHashStorage,
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
	skinHash, err := r.repo.GetSkinHash()
	if err != nil {
		return err
	}

	// 없어도 에러는 아님
	if skinHash != "" {
		r.configStorage.SaveConfigHash(consts.SKIN, skinHash)
	}

	// config 정보 init
	configHash, err := r.repo.GetConfigHash()
	if err != nil {
		return err
	}
	if configHash != "" {
		// 없어도 에러는 아님
		r.configStorage.SaveConfigHash(consts.CONFIG, configHash)
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

	serverConfig, err := r.configStorage.GetHash(consts.CONFIG)
	if err != nil {
		fmt.Println("config 에 대한 hash 정보를 찾을 수 없음.")
		configExist = false
	}

	// 있으면서 동일한
	if configExist && serverConfig == clientConfig {
		configSame = true
	}

	serverSkin, err := r.configStorage.GetHash(consts.SKIN)
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

func (r *commonUsecase) GetSkinImg(ctx context.Context, dto clDto.GetSkinImgRequest) (*os.File, error) {

	// skin hash 검증
	serverSkinHash, err := r.configStorage.GetHash(consts.SKIN)
	if err != nil {
		return nil, err
	}

	// 현재 서버 기준의 최신 skinHash와 클라이언트가 전달한 값이 다르면 처리하지 않음.
	if serverSkinHash != dto.SkinHash {
		return nil, consts.ErrSkinHashInvalid
	}

	// skin hash에 매핑된 파일 찾기
	filePath, err := r.configStorage.GetSkinFilePath(dto.SkinType)

	// 파일 존재 확인 정도는 usecase에서 할 수도 있음
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return file, nil
}
