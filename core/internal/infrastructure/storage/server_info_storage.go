package storage

import (
	"core/internal/domain/appValidation/entity"
	"errors"
	"fmt"
	"sync"
)

type ServerInfoStorage interface {
	SaveWorksCommonInfo(worksCode string, entity *entity.WorksCommonInfo) error
	GetWorksCommonInfo(worksCode string) *entity.WorksCommonInfo
}

type serverInfoStorage struct {
	mu              sync.RWMutex
	worksCommonInfo map[string]*entity.WorksCommonInfo
}

func NewServerInfoStorage() ServerInfoStorage {
	return &serverInfoStorage{
		worksCommonInfo: make(map[string]*entity.WorksCommonInfo),
	}
}

func (r *serverInfoStorage) SaveWorksCommonInfo(worksCode string, worksCommonInfo *entity.WorksCommonInfo) error {

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

func (r *serverInfoStorage) GetWorksCommonInfo(worksCode string) *entity.WorksCommonInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	worksCommonInfo, exists := r.worksCommonInfo[worksCode]
	if !exists {
		return nil
	}

	return worksCommonInfo
}
