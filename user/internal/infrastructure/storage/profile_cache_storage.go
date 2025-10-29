package storage

import (
	"log"
	"sync"
)

type profileCacheStorage struct {
	mu                  sync.RWMutex
	profileNameMap      map[string]string
	multiProfileNameMap map[string]string // 멀티 프로필 체크용 sha256(대상자id+내id)값으로 조회 (단, 멀티프로필을 1:1 로만 보여준다는 API라는 전제가 필요)
}

type ProfileCacheStorage interface {
	PutProfileName(userId string, fileName string)
	GetProfileName(userId string) string
	DeleteProfileName(userId string, fileName string)
}

func NewProfileCacheStorage() ProfileCacheStorage {
	return &profileCacheStorage{
		profileNameMap:      make(map[string]string),
		multiProfileNameMap: make(map[string]string),
	}
}

func (r *profileCacheStorage) PutProfileName(userId string, savedName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profileNameMap[userId] = savedName
	log.Printf("PutProfilePath : %s save success. \n", userId)
}

func (r *profileCacheStorage) GetProfileName(userId string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	savedPath, exists := r.profileNameMap[userId]
	if !exists {
		return ""
	} else {
		return savedPath
	}
}

func (r *profileCacheStorage) DeleteProfileName(userId string, fileName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	serverFileName := r.profileNameMap[userId]

	if serverFileName != "" && serverFileName == fileName {
		delete(r.profileNameMap, userId)
		log.Printf("[DeleteProfileName] profile name %s delete success. \n", fileName)
	}

}
