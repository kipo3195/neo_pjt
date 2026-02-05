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
