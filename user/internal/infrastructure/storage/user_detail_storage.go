package storage

import (
	"context"
	"log"
	"sync"
	"user/internal/domain/userDetail/entity"
)

type userDetailStorage struct {
	mu               sync.RWMutex
	updateVersionMap map[string]entity.UserUpdateVersionEntity
	userDetailMap    map[string]entity.UserDetailInfoEntity
}

type UserDetailStorage interface {
	GetUserInfoUpdateHash(ctx context.Context, reqUsers []entity.ReqUserEntity) (map[string]entity.UserUpdateVersionEntity, error)
	GetUserDetailInfo(userHash []string) []entity.UserDetailInfoEntity
	InitUserDetail(ctx context.Context, updateHash int64, initUsers []entity.UserDetailInfoEntity) error
}

func NewUserDetailStorage() UserDetailStorage {
	return &userDetailStorage{
		updateVersionMap: make(map[string]entity.UserUpdateVersionEntity),
		userDetailMap:    make(map[string]entity.UserDetailInfoEntity),
	}
}

func (r *userDetailStorage) GetUserInfoUpdateHash(ctx context.Context, reqUsers []entity.ReqUserEntity) (map[string]entity.UserUpdateVersionEntity, error) {
	// 1. 읽기 락 설정 (여러 고루틴이 동시에 읽을 수 있음)
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 2. 결과를 담을 맵 생성 (입력 크기만큼 용량 확보)
	result := make(map[string]entity.UserUpdateVersionEntity, len(reqUsers))

	// 3. 요청된 hash들만 찾아서 결과 맵에 담기
	for _, u := range reqUsers {
		if val, exists := r.updateVersionMap[u.UserHash]; exists {
			result[u.UserHash] = val
		}
	}
	return result, nil
}

func (r *userDetailStorage) GetUserDetailInfo(userHash []string) []entity.UserDetailInfoEntity {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.UserDetailInfoEntity, 0)

	for _, h := range userHash {
		temp, detailExists := r.userDetailMap[h]
		detailVersion, versionExists := r.updateVersionMap[h]
		if detailExists && versionExists {
			temp.DetailVersion = detailVersion.DetailVersion
			result = append(result, temp)
		} else {
			log.Printf("[GetUserDetailInfo] detailExists : %t, versionExists : %t\n", detailExists, versionExists)
		}
	}

	return result
}

// userDetailMap map[string]entity.UserDetailInfoEntity
func (r *userDetailStorage) InitUserDetail(ctx context.Context, detailVersion int64, entities []entity.UserDetailInfoEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println("[InitUserDetail] detailVersion : ", detailVersion)
	for _, i := range entities {

		temp := entity.UserDetailInfoEntity{
			Org:          i.Org,
			UserHash:     i.UserHash,
			UserPhoneNum: i.UserPhoneNum,
			UserEmail:    i.UserEmail,
		}
		r.userDetailMap[i.UserHash] = temp
		r.updateVersionMap[i.UserHash] = entity.UserUpdateVersionEntity{
			DetailVersion: detailVersion,
		}
	}

	return nil
}
