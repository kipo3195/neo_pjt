package storage

type userDetailStorage struct {
}

type UserDetailStorage interface {
}

func NewUserDetailStorage() UserDetailStorage {
	return &userDetailStorage{}
}
