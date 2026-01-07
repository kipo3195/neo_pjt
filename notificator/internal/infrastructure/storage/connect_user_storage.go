package storage

import (
	"log"
	"notificator/internal/domain/socketSender/entity"
	"sync"
)

type sendConnectionStorage struct {
	mu                 sync.RWMutex
	connectionMap      map[string]*entity.SendConnectionEntity
	connectionStateMap map[string]bool
}

type SendConnectionStorage interface {
	GetConnection(userHash string) *entity.SendConnectionEntity
	RemoveConnection(userHash string)
	PutConnection(userHash string, entity *entity.SendConnectionEntity)
}

func NewSendConnectionStorage() SendConnectionStorage {
	return &sendConnectionStorage{
		connectionMap:      make(map[string]*entity.SendConnectionEntity),
		connectionStateMap: make(map[string]bool),
	}
}

func (r *sendConnectionStorage) GetConnection(userHash string) *entity.SendConnectionEntity {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, exists := r.connectionMap[userHash]
	if !exists {
		return nil
	}

	return entity
}
func (r *sendConnectionStorage) RemoveConnection(userHash string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.connectionMap[userHash]
	if exists {
		delete(r.connectionMap, userHash)
	}
	log.Println("[RemoveConnection] userHash : ", userHash)
}

func (r *sendConnectionStorage) PutConnection(userHash string, entity *entity.SendConnectionEntity) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connectionMap[userHash] = entity
	log.Println("[PutConnection] userHash : ", userHash)
}
