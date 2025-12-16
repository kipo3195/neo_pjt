package repository

import "context"

type UserRepository interface {
	CheckUserRegist(ctx context.Context, id string) error
}
