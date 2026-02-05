package repository

import "context"

type ChatFileRepository interface {
	UpdateFileStatus(ctx context.Context, transactionId string) error
}
