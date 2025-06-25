package storage

import (
	"errors"
	"fmt"
	"sync"
)

type ConfigHashStorage interface {
	SaveConfigHash(kind string, hash string) error
	GetConfigHash(kind string) (string, error)
	DeleteConfigHash(kind string) error
}

type configHashStorage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewConfigHashStorage() ConfigHashStorage {
	return &configHashStorage{
		data: make(map[string]string),
	}
}

func (r *configHashStorage) SaveConfigHash(kind string, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[kind] = hash
	fmt.Printf("kind : %s hash : %s save success. \n", kind, hash)
	return nil
}

func (r *configHashStorage) GetConfigHash(kind string) (string, error) {

	hash, exists := r.data[kind]
	if !exists {
		return "", errors.New(kind + " is not exists")
	}

	return hash, nil
}

func (r *configHashStorage) DeleteConfigHash(kind string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[kind]; !exists {
		return errors.New(kind + "is not exists")
	}

	delete(r.data, kind)
	fmt.Printf("kind : %s delete success. \n", kind)
	return nil
}
