package storage

import (
	"errors"
	"sync"
)

type userStorage struct {
	mu               sync.RWMutex
	userChallengeMap map[string]string
}

type UserStorage interface {
	PutUserChallenge(userId string, challenge string) error
	GetUserChallenge(userId string) (string, error)
	DeleteUserChallenge(userId string) error
}

func NewUserStorage() UserStorage {
	return &userStorage{
		userChallengeMap: make(map[string]string),
	}
}

func (r *userStorage) PutUserChallenge(userId string, challenge string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.userChallengeMap[userId] = challenge

	return nil
}

func (r *userStorage) GetUserChallenge(userId string) (string, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	challenge, exists := r.userChallengeMap[userId]

	if !exists {
		return "", errors.New("challenge is not exists")
	}

	return challenge, nil
}

func (r *userStorage) DeleteUserChallenge(userId string) error {

	delete(r.userChallengeMap, userId)

	return nil
}
