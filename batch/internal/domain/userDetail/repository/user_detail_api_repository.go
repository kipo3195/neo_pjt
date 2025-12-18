package repository

import "context"

type UserDetailApiRepository interface {
	SendJsonToUser(ctx context.Context, fileName string, zipData []byte, orgCode string) error
}
