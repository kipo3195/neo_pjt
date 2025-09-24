package storage

import (
	"log"
	"sync"
)

type deviceStorage struct {
	mu                 sync.RWMutex
	deviceChallengeMap map[string]string //
	deviceAtTokenMap   map[string]string //
	deviceRtTokenMap   map[string]string //
}

type DeviceStorage interface {
	GetDeviceChallenge(Id string, uuid string) string
	PutDeviceChallenge(Id string, uuid string, chllange string)
	DeleteDeviceChallenge(Id string, uuid string)

	GetAccessToken(Id string, uuid string) string
	PutAccessToken(Id string, uuid string, at string)
	GetRefreshToken(Id string, uuid string) string
	PutRefreshToken(Id string, uuid string, rt string)
}

func NewDeviceStorage() DeviceStorage {
	return &deviceStorage{
		deviceChallengeMap: make(map[string]string),
		deviceAtTokenMap:   make(map[string]string),
		deviceRtTokenMap:   make(map[string]string),
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

func (r *deviceStorage) GetAccessToken(Id string, uuid string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	at, exists := r.deviceAtTokenMap[Id+":"+uuid]

	if !exists {
		return ""
	}
	return at
}

func (r *deviceStorage) PutAccessToken(Id string, uuid string, at string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deviceAtTokenMap[Id+":"+uuid] = at
	log.Println("PutAccessToken Id : ", Id)
}

func (r *deviceStorage) GetRefreshToken(Id string, uuid string) string {

	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, exists := r.deviceRtTokenMap[Id+":"+uuid]

	if !exists {
		return ""
	}

	return rt
}
func (r *deviceStorage) PutRefreshToken(Id string, uuid string, rt string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deviceRtTokenMap[Id+":"+uuid] = rt
	log.Println("PutRefreshToken Id : ", Id)
}
