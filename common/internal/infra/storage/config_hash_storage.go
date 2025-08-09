package storage

import (
	"common/internal/consts"
	"errors"
	"fmt"
	"sync"
)

type ConfigHashStorage interface {
	SaveConfigHash(kind string, hash string) error
	GetConfigHash() (string, error)
	DeleteConfigHash(kind string) error
}

type configHashStorage struct {
	mu         sync.RWMutex
	configInfo map[string]string // 해시 정보 저장용

}

func NewConfigHashStorage() ConfigHashStorage {
	return &configHashStorage{
		configInfo: make(map[string]string),
	}
}

func (r *configHashStorage) SaveConfigHash(config string, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.configInfo[config] = hash
	fmt.Printf("config : %s hash : %s save success. \n", config, hash)
	return nil
}

func (r *configHashStorage) GetConfigHash() (string, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	hash, exists := r.configInfo[consts.CONFIG]
	if !exists {
		return "", errors.New("configHash is not exists")
	}

	return hash, nil
}

func (r *configHashStorage) DeleteConfigHash(config string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.configInfo[config]; !exists {
		return errors.New(config + "is not exists")
	}

	delete(r.configInfo, config)
	fmt.Printf("config : %s delete success. \n", config)
	return nil
}
