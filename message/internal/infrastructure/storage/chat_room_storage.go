package storage

import (
	"message/internal/domain/chatRoom/entity"
	"sync"
)

type chatRoomStorage struct {
	mu            sync.RWMutex
	normalRoomMap map[string]entity.ChatRoomMemberEntity
}

type ChatRoomStorage interface {
}

func NewChatRoomStorage() ChatRoomStorage {
	return &chatRoomStorage{
		normalRoomMap: make(map[string]entity.ChatRoomMemberEntity),
	}
}
