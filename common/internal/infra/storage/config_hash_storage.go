package storage

import (
	"errors"
	"fmt"
	"sync"
)

type ConfigHashStorage interface {
	SaveConfigHash(kind string, hash string) error
	GetHash(kind string) (string, error)
	DeleteConfigHash(kind string) error
}

type configHashStorage struct {
	mu       sync.RWMutex
	hashInfo map[string]string // 해시 정보 저장용

}

func NewConfigHashStorage() ConfigHashStorage {
	return &configHashStorage{
		hashInfo: make(map[string]string),
	}
}

func (r *configHashStorage) SaveConfigHash(kind string, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hashInfo[kind] = hash
	fmt.Printf("kind : %s hash : %s save success. \n", kind, hash)
	return nil
}

func (r *configHashStorage) GetHash(kind string) (string, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	hash, exists := r.hashInfo[kind]
	if !exists {
		return "", errors.New(kind + " is not exists")
	}

	return hash, nil
}

func (r *configHashStorage) DeleteConfigHash(kind string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hashInfo[kind]; !exists {
		return errors.New(kind + "is not exists")
	}

	delete(r.hashInfo, kind)
	fmt.Printf("kind : %s delete success. \n", kind)
	return nil
}
