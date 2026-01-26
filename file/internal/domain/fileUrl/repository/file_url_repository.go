package repository

import (
	"context"
	"file/internal/domain/fileUrl/entity"
)

type FileUrlRepository interface {
	SaveCreateFileUrl(context context.Context, reqUserId string, transactionId string, en []entity.CreateFileUrlResultEntity) error
	GetFileId(ctx context.Context, en entity.FileUrlUploadEndEntity) ([]string, error)
}
