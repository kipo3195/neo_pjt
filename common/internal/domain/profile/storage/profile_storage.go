package storage

import (
	"context"
	"mime/multipart"
)

type ProfileStorage interface {
	Upload(ctx context.Context, file multipart.File, filename string) (string, error)
	GetProfileUrl(ctx context.Context, filename string) (string, error)
}
