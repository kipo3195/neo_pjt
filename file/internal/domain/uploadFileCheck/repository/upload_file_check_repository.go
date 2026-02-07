package repository

import (
	"context"
	"file/internal/domain/uploadFileCheck/entity"
)

type UploadFileCheckRepository interface {
	GetInvalidFile(ctx context.Context, checkDate string) ([]entity.InvalidFileEntity, error)
	UpdateInvalidFileState(ctx context.Context, fileIds []string) error
}
