package storage

type userAuthStorage struct {
}

type UserAuthStorage interface {
}

func NewUserAuthStorage() UserAuthStorage {
	return userAuthStorage{}
}
