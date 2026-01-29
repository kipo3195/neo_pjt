package cache

import (
	"context"
	"file/internal/domain/fileUrl/entity"
)

type FileUrlCache interface {
	PutFileUrlInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error
	GetFileUrlInfo(ctx context.Context, transactionId string) ([]entity.CreateFileUrlResultEntity, error)
	PutUploadEndFileInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error
}
