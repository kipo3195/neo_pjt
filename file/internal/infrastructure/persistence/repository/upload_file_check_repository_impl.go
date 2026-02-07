package repository

import (
	"context"
	"file/internal/domain/uploadFileCheck/entity"
	"file/internal/domain/uploadFileCheck/repository"
	"file/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type uploadFileCheckRepositoryImpl struct {
	db *gorm.DB
}

func NewUploadFileCheckRepository(db *gorm.DB) repository.UploadFileCheckRepository {
	return &uploadFileCheckRepositoryImpl{
		db: db,
	}
}

func (r *uploadFileCheckRepositoryImpl) GetInvalidFile(ctx context.Context, checkDate string) ([]entity.InvalidFileEntity, error) {

	var result []model.FileUploadUrlHistory

	err := r.db.WithContext(ctx).Raw(
		`select 
			file_id 
		from
			file_upload_url_history
		where 
			create_date like ? AND upload_flag ='N' AND error_flag ='N'`,
		checkDate+"%",
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	invalidFiles := make([]entity.InvalidFileEntity, 0)

	for _, value := range result {

		temp := entity.InvalidFileEntity{
			FileId: value.FileId,
		}
		invalidFiles = append(invalidFiles, temp)
	}

	return invalidFiles, nil
}

func (r *uploadFileCheckRepositoryImpl) UpdateInvalidFileState(ctx context.Context, invalidFileIds []string) error {

	return r.db.WithContext(ctx).Model(&model.FileUploadUrlHistory{}).Where(`file_id in ?`, invalidFileIds).Update("error_flag", "Y").Error
}
