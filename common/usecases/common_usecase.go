package usecases

import (
	"common/entities"
	"common/infra/storage"
	"common/internal/consts"
	"common/repositories"
)

type commonUsecase struct {
	repo          repositories.CommonRepository
	configStorage storage.ConfigHashStorage
}

type CommonUsecase interface {
	InitConfigHash() error
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
