package storage

import (
	"errors"
	"fmt"
	"sync"
)

type ConfigStorage interface {
	SaveConfigHash(kind string, hash string) error
	GetHash(kind string) (string, error)
	DeleteConfigHash(kind string) error
	GetSkinFilePath(skinType string) (string, error)
	SaveSkinFilePath(skinType string, filePath string) error
}

type configStorage struct {
	mu           sync.RWMutex
	hashInfo     map[string]string // 해시 정보 저장용
	skinFilePath map[string]string // 서버에서 접근할때 사용
}

func NewConfigHashStorage() ConfigStorage {
	return &configStorage{
		hashInfo:     make(map[string]string),
		skinFilePath: make(map[string]string),
	}
}

func (r *configStorage) SaveConfigHash(kind string, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hashInfo[kind] = hash
	fmt.Printf("kind : %s hash : %s save success. \n", kind, hash)
	return nil
}

func (r *configStorage) GetHash(kind string) (string, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	hash, exists := r.hashInfo[kind]
	if !exists {
		return "", errors.New(kind + " is not exists")
	}

	return hash, nil
}

func (r *configStorage) DeleteConfigHash(kind string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hashInfo[kind]; !exists {
		return errors.New(kind + "is not exists")
	}

	delete(r.hashInfo, kind)
	fmt.Printf("kind : %s delete success. \n", kind)
	return nil
}

func (r *configStorage) GetSkinFilePath(skinType string) (string, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath, exists := r.skinFilePath[skinType]

	if !exists {
		return "", errors.New(skinType + "is not exists")
	}
	return filePath, nil

}

func (r *configStorage) SaveSkinFilePath(skinType string, filePath string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	fmt.Printf("기존 skinType에 대한 파일 path skinType : %s, filePath : %s \n", skinType, r.skinFilePath[skinType])
	r.skinFilePath[skinType] = filePath
	fmt.Printf("변경된 skinType에 대한 파일 path skinType : %s, filePath : %s \n", skinType, r.skinFilePath[skinType])
	return nil

}
