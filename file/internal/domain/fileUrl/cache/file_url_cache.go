package cache

import (
	"context"
	"file/internal/domain/fileUrl/entity"
)

type FileUrlCache interface {
	PutFileUrlInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error
}
