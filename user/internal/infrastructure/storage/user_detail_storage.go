package storage

import (
	"sync"
	"user/internal/domain/userDetail/entity"
)

type userInfoServiceStorage struct {
	mu            sync.RWMutex
	updateHashMap map[string]string // userHash detailUpdateDate:profileUpdateDate 형식
	userDetailMap map[string]entity.UserDetailInfoEntity
}

type UserInfoServiceStorage interface {
}

func NewUserInfoServiceStorage() UserInfoServiceStorage {
	return &userInfoServiceStorage{
		updateHashMap: make(map[string]string),
		userDetailMap: make(map[string]entity.UserDetailInfoEntity),
	}
}
