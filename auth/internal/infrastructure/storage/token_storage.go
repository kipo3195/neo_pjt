package storage

import (
	"sync"
)

type tokenStorage struct {
	mu              sync.RWMutex
	refreshTokenMap map[string]string // 디바이스 인증 이후 rt 저장
	tokenExpMap     map[string]int    // 디바이스 인증 이후 at, rt (분) 만료시간 저장
}

type TokenStorage interface {
	PutTokenExp(t string, exp int)
	GetTokenExp(t string) int
}

func NewTokenStorage() TokenStorage {
	return &tokenStorage{
		refreshTokenMap: make(map[string]string),
		tokenExpMap:     make(map[string]int),
	}
}

func (r *tokenStorage) PutTokenExp(t string, exp int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokenExpMap[t] = exp
}

func (r *tokenStorage) GetTokenExp(t string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exp, exists := r.tokenExpMap[t]
	if !exists {
		return 0
	} else {
		return exp
	}
}
