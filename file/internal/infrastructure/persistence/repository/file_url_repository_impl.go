package repository

import (
	"context"
	"file/internal/consts"
	"file/internal/domain/fileUrl/cache"
	"file/internal/domain/fileUrl/entity"
	"file/internal/domain/fileUrl/repository"
	"file/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type fileUrlRepositoryImpl struct {
	db           *gorm.DB
	cacheStorage cache.FileUrlCache
}

func FileUrlMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.FileUploadUrlHistory{})
}

func NewFileUrlRepository(db *gorm.DB, cacheStorage cache.FileUrlCache) repository.FileUrlRepository {

	return &fileUrlRepositoryImpl{
		db:           db,
		cacheStorage: cacheStorage,
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
		}
	}

	err := r.db.WithContext(context).Create(fileUrlHistoryModels).Error
	if err != nil {
		return consts.ErrDB
	}

	r.cacheStorage.PutFileUrlInfo(context, transactionId, en)

	return nil
}

func (r *fileUrlRepositoryImpl) GetFileId(ctx context.Context, en entity.FileUrlUploadEndEntity) ([]entity.CreateFileUrlResultEntity, error) {

	fileInfo, err := r.cacheStorage.GetFileUrlInfo(ctx, en.TransactionId)

	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

func (r *fileUrlRepositoryImpl) UploadFlagUpdate(ctx context.Context, reqUserHash string, fileIds []string) error {

	err := r.db.WithContext(ctx).Model(&model.FileUploadUrlHistory{}).Where(`file_id in ?`, fileIds).Update("upload_flag", "Y").Error

	if err != nil {
		return err
	}

	return nil
}

func (r *fileUrlRepositoryImpl) PutUploadEndFileInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error {

	return r.cacheStorage.PutUploadEndFileInfo(ctx, transactionId, entity)
}
