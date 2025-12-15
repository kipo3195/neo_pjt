package storage

import (
	"log"
	"sync"
)

type orgStorage struct {
	mu         sync.RWMutex
	orgCode    []string
	orgCodeMap map[string]struct{}
}

type OrgStorage interface {
	GetOrgInfo() []string
	PutOrgInfo(org []string)
}

func NewOrgStorage() OrgStorage {
	return &orgStorage{
		orgCode:    make([]string, 0),
		orgCodeMap: make(map[string]struct{}),
	}
}

func (r *orgStorage) GetOrgInfo() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.orgCode
}

func (r *orgStorage) PutOrgInfo(org []string) {

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, value := range org {

		if _, exists := r.orgCodeMap[value]; !exists {
			r.orgCodeMap[value] = struct{}{}
			r.orgCode = append(r.orgCode, value)
		}

	}

	log.Println("[PutOrgInfo] orgCode info :", r.orgCode)
}
