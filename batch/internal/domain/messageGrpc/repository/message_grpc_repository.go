package repository

import (
	"context"
)

type MessageGrpcRepository interface {
	GetSendFileInfo(ctx context.Context, checkDate string, fileIds []string) (map[string]string, error)
}
