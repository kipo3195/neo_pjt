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

func (r *fileUrlApiRepositoryImpl) CreateFileUrl(ctx context.Context, en entity.CreateFileUrlEntity) ([]entity.CreateFileUrlResultEntity, error) {

	result := make([]entity.CreateFileUrlResultEntity, 0)

	for key := range en.FileInfoMap {

		fileInfoEntity := en.FileInfoMap[key]

		temp := entity.CreateFileUrlResultEntity{
			FileId:     key,
			FileName:   fileInfoEntity.FileName,
			CreatedUrl: key + " url",
		}

		result = append(result, temp)
	}

	return result, nil
}
