package repository

import (
	"context"
	"file/internal/domain/fileUrl/entity"
	"file/internal/domain/fileUrl/repository"
)

type fileUrlApiRepositoryImpl struct {
}

func NewFileUrlApiRepository() repository.FileUrlApiRepository {
	return &fileUrlApiRepositoryImpl{}
}

func (r *fileUrlApiRepositoryImpl) CreateFileUrl(ctx context.Context, entity entity.CreateFileUrlEntity) ([]entity.CreateFileUrlResultEntity, error) {

	return nil, nil
}
