package storage

import (
	"errors"
	"fmt"
	"sync"
)

type skinStorage struct {
	mu           sync.RWMutex
	skinFilePath map[string]string // 서버에서 접근할때 사용
}

type SkinStorage interface {
	GetSkinFilePath(skinType string) (string, error)
	SaveSkinFilePath(skinType string, filePath string) error
}

func NewSkinStorage() SkinStorage {
	return &skinStorage{
		skinFilePath: make(map[string]string),
	}
}

func (r *skinStorage) GetSkinFilePath(skinType string) (string, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath, exists := r.skinFilePath[skinType]

	if !exists {
		return "", errors.New(skinType + "is not exists")
	}
	return filePath, nil

}

func (r *skinStorage) SaveSkinFilePath(skinType string, filePath string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	fmt.Printf("기존 skinType에 대한 파일 path skinType : %s, filePath : %s \n", skinType, r.skinFilePath[skinType])
	r.skinFilePath[skinType] = filePath
	fmt.Printf("변경된 skinType에 대한 파일 path skinType : %s, filePath : %s \n", skinType, r.skinFilePath[skinType])
	return nil

}
