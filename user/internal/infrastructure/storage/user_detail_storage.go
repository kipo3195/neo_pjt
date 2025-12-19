package storage

import (
	"context"
	"sync"
	"user/internal/domain/userDetail/entity"
)

type userInfoServiceStorage struct {
	mu            sync.RWMutex
	updateHashMap map[string]entity.UserUpdateHashEntity
	userDetailMap map[string]entity.UserDetailInfoEntity
}

type UserInfoServiceStorage interface {
	GetUserInfoUpdateHash(ctx context.Context, reqUsers []entity.ReqUserEntity) (map[string]entity.UserUpdateHashEntity, error)
	GetUserDetailInfo(userHash []string) []entity.UserDetailInfoEntity
}

func NewUserInfoServiceStorage() UserInfoServiceStorage {
	return &userInfoServiceStorage{
		updateHashMap: make(map[string]entity.UserUpdateHashEntity),
		userDetailMap: make(map[string]entity.UserDetailInfoEntity),
	}
}

func (r *userInfoServiceStorage) GetUserInfoUpdateHash(ctx context.Context, reqUsers []entity.ReqUserEntity) (map[string]entity.UserUpdateHashEntity, error) {
	// 1. 읽기 락 설정 (여러 고루틴이 동시에 읽을 수 있음)
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 2. 결과를 담을 맵 생성 (입력 크기만큼 용량 확보)
	result := make(map[string]entity.UserUpdateHashEntity, len(reqUsers))

	// 3. 요청된 hash들만 찾아서 결과 맵에 담기
	for _, u := range reqUsers {
		if val, exists := r.updateHashMap[u.UserHash]; exists {
			result[u.UserHash] = val
		}
	}
	return result, nil
}

func (r *userInfoServiceStorage) GetUserDetailInfo(userHash []string) []entity.UserDetailInfoEntity {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.UserDetailInfoEntity, 0)

	for _, h := range userHash {
		temp, exists := r.userDetailMap[h]
		if exists {
			result = append(result, temp)
		}
	}

	return result
}
