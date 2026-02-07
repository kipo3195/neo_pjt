package usecase

import (
	"context"
	"file/internal/domain/uploadFileCheck/repository"
)

type uploadFileCheckUsecase struct {
	repository repository.UploadFileCheckRepository
}

type UploadFileCheckUsecase interface {
	UploadFileCheck(ctx context.Context, checkDate string) error
}

func NewUploadFileCheckUsecase(repository repository.UploadFileCheckRepository) UploadFileCheckUsecase {
	return &uploadFileCheckUsecase{
		repository: repository,
	}
}

func (r *uploadFileCheckUsecase) UploadFileCheck(ctx context.Context, checkDate string) error {
	return nil
}
