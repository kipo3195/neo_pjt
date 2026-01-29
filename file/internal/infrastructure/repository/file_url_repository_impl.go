package repository

import (
	"context"
	"file/internal/consts"
	"file/internal/domain/fileUrl/cache"
	"file/internal/domain/fileUrl/entity"
	"file/internal/domain/fileUrl/repository"
	"file/internal/infrastructure/model"

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
			UploadUrl:     file.CreatedUrl,
		}
	}

	err := r.db.WithContext(context).Create(fileUrlHistoryModels).Error
	if err != nil {
		return consts.ErrDB
	}

	r.cacheStorage.PutFileUrlInfo(context, transactionId, en)

	return nil
}

func (r *fileUrlRepositoryImpl) GetFileId(ctx context.Context, en entity.FileUrlUploadEndEntity) ([]string, error) {

	var m []entity.FileUploadUrlHistoryEntity

	err := r.db.WithContext(ctx).Raw(
		`select file_id, upload_flag, error_flag from file_upload_url_history where t_id = ? and req_user_hash = ?`, en.TransactionId, en.ReqUserHash,
	).Scan(&m).Error

	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	for _, v := range m {

		if v.UploadFlag == "Y" && v.ErrorFlag == "Y" {
			return nil, consts.ErrULIDGeneratorError
		}

		result = append(result, v.FileId)
	}

	return result, nil
}
