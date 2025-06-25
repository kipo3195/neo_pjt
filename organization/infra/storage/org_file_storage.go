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

func (m *orgFileStorage) SaveOrgFile(orgCode string, content []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// content 복사해서 저장 (안전하게)
	copied := make([]byte, len(content))
	copy(copied, content)
	m.data[orgCode] = copied

	return nil
}

func (m *orgFileStorage) GetOrgFile(orgCode string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.data[orgCode]
	if !exists {
		return nil, errors.New("zip file not found for org: " + orgCode)
	}

	// 안전하게 복사본 반환
	copied := make([]byte, len(data))
	copy(copied, data)

	return copied, nil
}

func (m *orgFileStorage) DeleteOrgFile(orgCode string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[orgCode]; !exists {
		return errors.New("file not found")
	}

	delete(m.data, orgCode)
	return nil
}
