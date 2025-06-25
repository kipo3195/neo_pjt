package storage

import (
	"errors"
	"sync"
)

type OrgFileStorage interface {
	SaveOrgFile(orgID string, content []byte) error
	GetOrgFile(orgID string) ([]byte, error)
	DeleteOrgFile(orgID string) error
}

type orgFileStorage struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewOrgFileStorage() OrgFileStorage {
	return &orgFileStorage{
		data: make(map[string][]byte),
	}
}

func (r *orgFileStorage) SaveOrgFile(orgCode string, content []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// content 복사해서 저장 (안전하게)
	copied := make([]byte, len(content))
	copy(copied, content)
	r.data[orgCode] = copied

	return nil
}

func (r *orgFileStorage) GetOrgFile(orgCode string) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, exists := r.data[orgCode]
	if !exists {
		return nil, errors.New("zip file not found for org: " + orgCode)
	}

	// 안전하게 복사본 반환
	copied := make([]byte, len(data))
	copy(copied, data)

	return copied, nil
}

func (r *orgFileStorage) DeleteOrgFile(orgCode string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[orgCode]; !exists {
		return errors.New("file not found")
	}

	delete(r.data, orgCode)
	return nil
}
