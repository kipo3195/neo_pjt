package storage

import (
	"log"
	"sync"
)

type serviceUserStorage struct {
	mu             sync.RWMutex
	serviceUserMap map[string]string // id : hash
}

type ServiceUserStorage interface {
	GetUserHash(Id string) string
	PutUserhash(Id string, hash string)
}

func NewServiceUserStorage() ServiceUserStorage {
	return &serviceUserStorage{
		serviceUserMap: make(map[string]string),
	}
}

func (r *serviceUserStorage) GetUserHash(Id string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hash, exists := r.serviceUserMap[Id]
	if !exists {
		log.Printf("GetUserHash %s hash is null.", Id)
		return ""
	}
	return hash
}

func (r *serviceUserStorage) PutUserhash(Id string, hash string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.serviceUserMap[Id] = hash
}
