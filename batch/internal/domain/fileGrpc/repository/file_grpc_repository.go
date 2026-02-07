package repository

import "context"

type FileGrpcRepository interface {
	CheckUploadFile(ctx context.Context, checkDate string) error
}
