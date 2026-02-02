package repository

import "context"

type ChatApiRepository interface {
	NotifySendChatFile(ctx context.Context, transactionId string) error
}
