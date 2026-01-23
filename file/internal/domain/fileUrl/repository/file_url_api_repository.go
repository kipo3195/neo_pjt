package repository

import (
	"context"
	"file/internal/domain/fileUrl/entity"
)

type FileUrlApiRepository interface {
	CreateFileUrl(ctx context.Context, entity entity.CreateFileUrlEntity) ([]entity.CreateFileUrlResultEntity, error)
}
