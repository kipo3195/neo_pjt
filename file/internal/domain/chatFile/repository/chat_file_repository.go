package repository

import "context"

type ChatFileRepository interface {
	UpdateFileStatus(ctx context.Context, transactionId string) error
	GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error)
	SendFlagUpdate(ctx context.Context, sendedFileIds []string) error
}
