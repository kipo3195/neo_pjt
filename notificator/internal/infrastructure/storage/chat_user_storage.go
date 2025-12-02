package storage

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type chatUserStorage struct {
	mu                 sync.RWMutex
	chatUserConnectMap map[string]*websocket.Conn
}

type ChatUserStorage interface {
	GetChatConnect(userHash string) *websocket.Conn
	RemoveChatConnect(userHash string)
	PutChatConnect(userHash string, conn *websocket.Conn, c chan []byte)
}

func NewChatUserStorage() ChatUserStorage {
	return &chatUserStorage{
		chatUserConnectMap: make(map[string]*websocket.Conn),
	}
}

func (r *chatUserStorage) GetChatConnect(userHash string) *websocket.Conn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.chatUserConnectMap[userHash]
}
func (r *chatUserStorage) RemoveChatConnect(userHash string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	connect := r.chatUserConnectMap[userHash]

	if connect != nil {
		delete(r.chatUserConnectMap, userHash)
	}
	log.Println("[RemoveChatConnect] userHash : ", userHash)
}
func (r *chatUserStorage) PutChatConnect(userHash string, conn *websocket.Conn, c chan []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.chatUserConnectMap[userHash] = conn
	log.Println("[PutChatConnect] userHash : ", userHash)

}
