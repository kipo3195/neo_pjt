package client

import (
	"common/internal/consts"
	"common/internal/domains/skin/dto/client/requestDTO"
	repositories "common/internal/domains/skin/repositories/client"
	"common/internal/infra/storage"
	"context"
	"fmt"
	"log"
	"os"
)

type skinUsecase struct {
	repository    repositories.SkinRepository
	skinStorage   storage.SkinStorage
	configStorage storage.ConfigHashStorage
}

type SkinUsecase interface {
	GetSkinImg(ctx context.Context, dto requestDTO.GetSkinImgRequest) (*os.File, error)
	CheckSkin(skinHash string) (bool, error)
}

func NewSkinUsecase(repository repositories.SkinRepository, skinStorage storage.SkinStorage, configStorage storage.ConfigHashStorage) SkinUsecase {
	return &skinUsecase{
		repository:    repository,
		skinStorage:   skinStorage,
		configStorage: configStorage,
	}
}

func (r *skinUsecase) GetSkinImg(ctx context.Context, dto requestDTO.GetSkinImgRequest) (*os.File, error) {

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
	filePath, err := r.skinStorage.GetSkinFilePath(dto.SkinType)

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

func (r *skinUsecase) CheckSkin(skinHash string) (bool, error) {

	skinHash, err := r.skinStorage.GetSkinHash()
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
