package storage

import (
	"sync"
)

type authTokenStorage struct {
	mu            sync.RWMutex
	atTokenMap    map[string]string //
	rtTokenMap    map[string]string //
	rtTokenExpMap map[string]string //
}

type AuthTokenStorage interface {
	GetAccessToken(Id string, uuid string) string
	PutAccessToken(Id string, uuid string, at string)
	GetRefreshToken(Id string, uuid string) string
	PutRefreshToken(Id string, uuid string, rt string)
	GetRefreshTokenExp(Id string, uuid string) string
	PutRefreshTokenExp(Id string, uuid string, rtExp string)
}

func NewAuthTokenStorage() AuthTokenStorage {
	return &authTokenStorage{
		atTokenMap:    make(map[string]string),
		rtTokenMap:    make(map[string]string),
		rtTokenExpMap: make(map[string]string),
	}
}

func (r *authTokenStorage) GetAccessToken(Id string, uuid string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	at, exists := r.atTokenMap[Id+":"+uuid]

	if !exists {
		return ""
	}
	return at
}

func (r *authTokenStorage) PutAccessToken(Id string, uuid string, at string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.atTokenMap[Id+":"+uuid] = at
}

func (r *authTokenStorage) GetRefreshToken(Id string, uuid string) string {

	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, exists := r.rtTokenMap[Id+":"+uuid]

	if !exists {
		return ""
	}

	return rt
}

func (r *authTokenStorage) PutRefreshToken(Id string, uuid string, rt string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rtTokenMap[Id+":"+uuid] = rt
}

func (r *authTokenStorage) GetRefreshTokenExp(Id string, uuid string) string {

	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, exists := r.rtTokenExpMap[Id+":"+uuid]

	if !exists {
		return ""
	}

	return rt
}

func (r *authTokenStorage) PutRefreshTokenExp(Id string, uuid string, rtExp string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rtTokenExpMap[Id+":"+uuid] = rtExp
}
