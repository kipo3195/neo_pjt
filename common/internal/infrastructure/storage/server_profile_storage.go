package storage

import (
	"context"
	"mime/multipart"
)

// profile 도메인의 storage에 정의된 행동 계약(Contract)만 정의
type ServerProfileStorage struct {
	ServerUrl string
}

func NewServerProfileStorage(serverUrl string) *ServerProfileStorage {
	return &ServerProfileStorage{
		ServerUrl: serverUrl,
	}
}

func (s *ServerProfileStorage) Upload(ctx context.Context, file multipart.File, filename string) (string, error) {
	return "", nil
}

func (s *ServerProfileStorage) GetProfileUrl(ctx context.Context, filename string) (string, error) {

	return "", nil
}
