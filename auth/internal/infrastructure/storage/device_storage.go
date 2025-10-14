package storage

import (
	"log"
	"sync"
)

type deviceStorage struct {
	mu                 sync.RWMutex
	deviceChallengeMap map[string]string // 디바이스 인증을 위한 challenge 발급 저장
}

type DeviceStorage interface {
	GetDeviceChallenge(Id string, uuid string) string
	PutDeviceChallenge(Id string, uuid string, chllange string)
	DeleteDeviceChallenge(Id string, uuid string)
}

func NewDeviceStorage() DeviceStorage {
	return &deviceStorage{
		deviceChallengeMap: make(map[string]string),
	}
}

func (r *deviceStorage) GetDeviceChallenge(Id string, uuid string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	challenge, exists := r.deviceChallengeMap[Id+":"+uuid]

	if !exists {
		return ""
	}

	return challenge
}

func (r *deviceStorage) PutDeviceChallenge(Id string, uuid string, challenge string) {

	r.mu.Lock()
	defer r.mu.Unlock()
	r.deviceChallengeMap[Id+":"+uuid] = challenge
	log.Println("PutDeviceChallenge Id : ", Id)

}
func (r *deviceStorage) DeleteDeviceChallenge(Id string, uuid string) {
	delete(r.deviceChallengeMap, Id+":"+uuid)
}
