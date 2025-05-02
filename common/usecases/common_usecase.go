package usecases

import (
	"common/dto"
	"common/entities"
	"common/repositories"
	"fmt"
	"os"
	"path/filepath"
)

type commonUsecase struct {
	repo repositories.CommonRepository
}

type CommonUsecase interface {
	GetConfig(dto.ConfigRequest) (*entities.Config, error)
}

func NewCommonUsecase(repo repositories.CommonRepository) CommonUsecase {
	return &commonUsecase{repo: repo}
}

func (u *commonUsecase) GetConfig(req dto.ConfigRequest) (*entities.Config, error) {
	// 대부분의 시스템에서는 단일 파일 다운로드 시 다음과 같은 패턴을 따릅니다:
	// Content-Type: application/octet-stream 또는 해당 파일 타입 (예: application/json, text/plain, application/x-yaml 등)

	// 설정으로 주입해야 하나?
	fileName := req.FileName
	filePath := filepath.Join("./config", fileName)

	// 파일 읽기
	content, err := os.ReadFile(filePath)
	fmt.Println("읽은 파일 ", content)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return &entities.Config{
		FileName: fileName,
		Content:  content,
	}, nil
}
