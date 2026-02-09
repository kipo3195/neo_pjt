package repository

import (
	"context"
	"file/internal/domain/chatFile/repository"
	"file/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type chatFileRespositoryImpl struct {
	db *gorm.DB
}

func NewChatFileRepository(db *gorm.DB) repository.ChatFileRepository {
	return &chatFileRespositoryImpl{
		db: db,
	}
}

func (r *chatFileRespositoryImpl) UpdateFileStatus(ctx context.Context, transactionId string) error {

	err := r.db.WithContext(ctx).Model(&model.FileUploadUrlHistory{}).Where(`t_id = ?`, transactionId).Update("send_flag", "Y").Error

	if err != nil {
		return err
	}
	return nil
}

func (r *chatFileRespositoryImpl) GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error) {

	var result []model.FileUploadUrlHistory

	err := r.db.WithContext(ctx).Raw(
		`select 
			file_id 
		from 
			file_upload_url_history 
		where 
			create_date like ? AND send_flag ='N' and error_flag ='N'`,
		yesterday+"%",
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	invalidFileIds := make([]string, 0)

	for _, value := range result {
		invalidFileIds = append(invalidFileIds, value.FileId)
	}

	return invalidFileIds, nil
}

func (r *chatFileRespositoryImpl) SendFlagUpdate(ctx context.Context, sendedFileIds []string) error {

	err := r.db.WithContext(ctx).Model(&model.FileUploadUrlHistory{}).Where(`file_id in ?`, sendedFileIds).Update("send_flag", "Y").Error

	if err != nil {
		return err
	}

	return nil
}
