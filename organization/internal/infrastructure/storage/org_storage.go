package storage

import (
	"log"
	"sync"
)

type OrgStorage interface {
	PutWorksOrgCode(orgCode []string)
	GetWorksOrgCode() []string
}

type orgStorage struct {
	mu              sync.RWMutex
	worksOrgCode    []string
	worksOrgCodeMap map[string]struct{}
}

func NewOrgStorage() OrgStorage {
	return &orgStorage{
		worksOrgCode:    make([]string, 0),
		worksOrgCodeMap: make(map[string]struct{}),
	}
}

func (r *orgStorage) PutWorksOrgCode(orgCode []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, code := range orgCode {

		if _, exists := r.worksOrgCodeMap[code]; !exists {
			r.worksOrgCodeMap[code] = struct{}{}
			r.worksOrgCode = append(r.worksOrgCode, code)
		}
	}
	log.Println("[PutWorksOrgCode] orgCode : ", r.worksOrgCode)
}

func (r *orgStorage) GetWorksOrgCode() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.worksOrgCode
}
