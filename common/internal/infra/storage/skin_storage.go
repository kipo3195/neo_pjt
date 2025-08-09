package storage

import (
	"common/internal/consts"
	"errors"
	"fmt"
	"sync"
)

type skinStorage struct {
	mu           sync.RWMutex
	skinFilePath map[string]string // 서버에서 접근할때 사용
	hashInfo     map[string]string // 해시 정보 저장용
}

type SkinStorage interface {
	GetSkinFilePath(skinType string) (string, error)
	SaveSkinFilePath(skinType string, filePath string) error
	GetSkinHash() (string, error)
}

func NewSkinStorage() SkinStorage {
	return &skinStorage{
		skinFilePath: make(map[string]string),
		hashInfo:     make(map[string]string), // 해시 정보 저장용
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

func (r *skinStorage) GetSkinHash() (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	serverHash, exists := r.hashInfo[consts.SKIN]
	if !exists {
		return "", errors.New("skin hash is not exists")
	}

	return serverHash, nil
}
