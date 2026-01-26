package repository

import (
	"context"
	"file/internal/domain/fileUrl/entity"
	"file/internal/domain/fileUrl/repository"
	"file/internal/infrastructure/model"

	"gorm.io/gorm"
)

type fileUrlRepositoryImpl struct {
	db *gorm.DB
}

func FileUrlMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.FileUploadUrlHistory{})
}

func NewFileUrlRepository(db *gorm.DB) repository.FileUrlRepository {

	return &fileUrlRepositoryImpl{
		db: db,
	}
}

func (r *fileUrlRepositoryImpl) SaveCreateFileUrl(context context.Context, reqUserId string, transactionId string, en []entity.CreateFileUrlResultEntity) error {

	fileUrlHistoryModels := make([]model.FileUploadUrlHistory, len(en))

	for i, file := range en {
		fileUrlHistoryModels[i] = model.FileUploadUrlHistory{
			FileId:        file.FileId,
			ReqUserHash:   reqUserId,
			TransactionId: transactionId,
			FileName:      file.FileName,
			UploadUrl:     file.CreatedUrl,
		}
	}

	return r.db.WithContext(context).Create(fileUrlHistoryModels).Error
}
