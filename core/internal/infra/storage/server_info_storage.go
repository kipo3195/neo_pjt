package storage

import (
	"core/internal/domains/appValidation/entities"
	"errors"
	"fmt"
	"sync"
)

type ServerInfoStorage interface {
	SaveWorksCommonInfo(worksCode string, entity *entities.WorksCommonInfo) error
	GetWorksCommonInfo(worksCode string) *entities.WorksCommonInfo
}

type serverInfoStorage struct {
	mu              sync.RWMutex
	worksCommonInfo map[string]*entities.WorksCommonInfo
}

func NewServerInfoStorage() ServerInfoStorage {
	return &serverInfoStorage{
		worksCommonInfo: make(map[string]*entities.WorksCommonInfo),
	}
}

func (r *serverInfoStorage) SaveWorksCommonInfo(worksCode string, worksCommonInfo *entities.WorksCommonInfo) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if worksCode == "" {
		return errors.New(worksCode + " is invalid")
	}

	if worksCommonInfo == nil {
		return errors.New("worksCommonInfo is null")
	}

	r.worksCommonInfo[worksCode] = worksCommonInfo
	fmt.Printf("[SaveWorksCommonInfo] %s worksCommonInfo save success !\n", worksCode)

	return nil
}

func (r *serverInfoStorage) GetWorksCommonInfo(worksCode string) *entities.WorksCommonInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	worksCommonInfo, exists := r.worksCommonInfo[worksCode]
	if !exists {
		return nil
	}

	return worksCommonInfo
}
