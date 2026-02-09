package repository

import "context"

type FileGrpcRepository interface {
	CheckUploadFile(ctx context.Context, checkDate string) error
	GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error)
	ClearFileStorage(ctx context.Context, clearFileId []string, sendedFileId []string) error
}
