package storage

import (
	"log"
	"sync"
)

type userAuthStorage struct {
	mu               sync.RWMutex
	userChallengeMap map[string]string //
}

type UserAuthStorage interface {
	GetUserAuthChallenge(Id string) string
	PutUserAuthChallenge(Id string, challenge string)
	DeleteUserAuthChallenge(Id string)
}

func NewUserAuthStorage() UserAuthStorage {
	return &userAuthStorage{
		userChallengeMap: make(map[string]string),
	}
}

func (r *userAuthStorage) GetUserAuthChallenge(Id string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	challenge, exists := r.userChallengeMap[Id]

	if !exists {
		log.Printf("GetUserAuthChallenge %s challenge is null.", Id)
		return ""
	}
	return challenge
}

func (r *userAuthStorage) PutUserAuthChallenge(Id string, challenge string) {

	r.mu.Lock()
	defer r.mu.Unlock()
	r.userChallengeMap[Id] = challenge
	log.Println("PutUserAuthChallenge Id : ", Id)

}

func (r *userAuthStorage) DeleteUserAuthChallenge(Id string) {
	delete(r.userChallengeMap, Id)
}
