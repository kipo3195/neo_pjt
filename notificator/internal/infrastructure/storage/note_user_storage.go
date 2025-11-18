package storage

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type noteUserStorage struct {
	mu                 sync.RWMutex
	noteUserConnectMap map[string]*websocket.Conn
}

type NoteUserStorage interface {
	GetNoteConnect(userHash string) *websocket.Conn
	RemoveNoteConnect(userHash string)
	PutNoteConnect(userHash string, conn *websocket.Conn)
}

func NewNoteUserStorage() NoteUserStorage {
	return &noteUserStorage{
		noteUserConnectMap: make(map[string]*websocket.Conn),
	}
}

func (r *noteUserStorage) GetNoteConnect(userHash string) *websocket.Conn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.noteUserConnectMap[userHash]
}
func (r *noteUserStorage) RemoveNoteConnect(userHash string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	connect := r.noteUserConnectMap[userHash]

	if connect != nil {
		delete(r.noteUserConnectMap, userHash)
	}
	log.Println("[RemoveNoteConnect] userHash : ", userHash)
}
func (r *noteUserStorage) PutNoteConnect(userHash string, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.noteUserConnectMap[userHash] = conn
	log.Println("[PutNoteConnect] userHash : ", userHash)

}
