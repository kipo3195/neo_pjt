package repository

import (
	"context"
	"file/internal/domain/fileUrl/entity"
)

type FileUrlRepository interface {
	SaveCreateFileUrl(context context.Context, reqUserHash string, transactionId string, en []entity.CreateFileUrlResultEntity) error
	GetFileId(ctx context.Context, en entity.FileUrlUploadEndEntity) ([]entity.CreateFileUrlResultEntity, error)
	UploadFlagUpdate(ctx context.Context, reqUserHash string, fileIds []string) error
	PutUploadEndFileInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error
}
