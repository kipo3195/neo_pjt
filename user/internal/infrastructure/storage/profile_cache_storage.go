package storage

import (
	"context"
	"log"
	"sync"
	"user/internal/domain/profile/entity"
)

type profileCacheStorage struct {
	mu                      sync.RWMutex
	profileUpdateVersionMap map[string]entity.ProfileUpdateVersionEntity
	profileNameMap          map[string]string
	multiProfileNameMap     map[string]string // 멀티 프로필 체크용 sha256(대상자id+내id)값으로 조회 (단, 멀티프로필을 1:1 로만 보여준다는 API라는 전제가 필요)
}

type ProfileCacheStorage interface {
	PutProfileName(userHash string, fileName string)
	GetProfileName(userHash string) string
	DeleteProfileName(userHash string, fileName string)
	InitProfile(ctx context.Context, profileVersion int64, entity []entity.ProfileInfoEntity)
	PutProfileVersion(userHash string, profileVersion int64)
	GetUserProfileUpdateVersion(ctx context.Context, entity []entity.ReqUserEntity) (map[string]entity.ProfileUpdateVersionEntity, error)
	GetProfileInfo(userHash []string) []entity.ProfileInfoEntity
}

func NewProfileCacheStorage() ProfileCacheStorage {
	return &profileCacheStorage{
		profileNameMap:          make(map[string]string),
		multiProfileNameMap:     make(map[string]string),
		profileUpdateVersionMap: make(map[string]entity.ProfileUpdateVersionEntity),
	}
}

func (r *profileCacheStorage) PutProfileName(userHash string, savedName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profileNameMap[userHash] = savedName
	log.Printf("PutProfilePath : %s save success. \n", userHash)
}

func (r *profileCacheStorage) GetProfileName(userHash string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	savedPath, exists := r.profileNameMap[userHash]
	if !exists {
		return ""
	} else {
		return savedPath
	}
}

func (r *profileCacheStorage) DeleteProfileName(userHash string, fileName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	serverFileName := r.profileNameMap[userHash]

	if serverFileName != "" && serverFileName == fileName {
		delete(r.profileNameMap, userHash)
		log.Printf("[DeleteProfileName] profile name %s delete success. \n", fileName)
	}

}

func (r *profileCacheStorage) InitProfile(ctx context.Context, profileVersion int64, entities []entity.ProfileInfoEntity) {

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, i := range entities {
		r.profileNameMap[i.UserHash] = i.UserHash
		r.profileUpdateVersionMap[i.UserHash] = entity.ProfileUpdateVersionEntity{
			ProfileVersion: profileVersion,
		}
	}

}

func (r *profileCacheStorage) PutProfileVersion(userHash string, profileVersion int64) {

	r.mu.Lock()
	defer r.mu.Unlock()
	r.profileUpdateVersionMap[userHash] = entity.ProfileUpdateVersionEntity{
		ProfileVersion: profileVersion,
	}

	log.Println("[PutProfileVersion] map result : ", r.profileUpdateVersionMap)
}

func (r *profileCacheStorage) GetUserProfileUpdateVersion(ctx context.Context, reqUsers []entity.ReqUserEntity) (map[string]entity.ProfileUpdateVersionEntity, error) {

	// 1. 읽기 락 설정 (여러 고루틴이 동시에 읽을 수 있음)
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 2. 결과를 담을 맵 생성 (입력 크기만큼 용량 확보)
	result := make(map[string]entity.ProfileUpdateVersionEntity, len(reqUsers))

	// 3. 요청된 hash들만 찾아서 결과 맵에 담기
	for _, u := range reqUsers {
		if val, exists := r.profileUpdateVersionMap[u.UserHash]; exists {
			result[u.UserHash] = val
		}
	}
	return result, nil
}

func (r *profileCacheStorage) GetProfileInfo(userHash []string) []entity.ProfileInfoEntity {

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.ProfileInfoEntity, 0)

	for _, h := range userHash {
		profileName, profileExists := r.profileNameMap[h]
		profileVersion, versionExists := r.profileUpdateVersionMap[h]
		if profileExists && versionExists {

			temp := entity.ProfileInfoEntity{
				UserHash:       h,
				SaveName:       profileName,
				ProfileVersion: profileVersion.ProfileVersion,
			}

			result = append(result, temp)
		} else {
			log.Printf("[GetProfileInfo] profileName : %t, versionExists : %t\n", profileExists, versionExists)
		}
	}

	return result
}
