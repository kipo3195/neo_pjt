package storage

import (
	"log"
	"notificator/internal/domain/chat/entity"
	"sync"

	"github.com/gorilla/websocket"
)

type chatUserStorage struct {
	mu                 sync.RWMutex
	chatUserConnectMap map[string]*websocket.Conn
	chatRoomMemberMap  map[string][]string
	chatRoomMemberSet  map[string]map[string]struct{}
}

type ChatUserStorage interface {
	GetChatConnect(userHash string) *websocket.Conn
	RemoveChatConnect(userHash string)
	PutChatConnect(userHash string, conn *websocket.Conn, c chan []byte)
	GetChatRoomMember(roomKey string) []string
	PutChatRoomMember(roomKey string, member []entity.CreateChatRoomMemberEntity)
}

func NewChatUserStorage() ChatUserStorage {
	return &chatUserStorage{
		chatUserConnectMap: make(map[string]*websocket.Conn),
		chatRoomMemberMap:  make(map[string][]string),
		chatRoomMemberSet:  make(map[string]map[string]struct{}),
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
func (r *chatUserStorage) GetChatRoomMember(roomKey string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// mutex 는 지연이 있을 수 도 있지 않을까? 그러니까 sync 자료구조로?
	log.Println("[GetChatRoomMember]", roomKey)

	member, exists := r.chatRoomMemberMap[roomKey]
	if !exists {
		log.Println("[GetChatRoomMember] roomKey is not regist.")
		return nil
	}

	// slice 복사 (중요)
	result := make([]string, len(member))
	copy(result, member)

	return result
}

func (r *chatUserStorage) PutChatRoomMember(roomKey string, member []entity.CreateChatRoomMemberEntity) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// set으로 중복 제거
	unique := make(map[string]struct{}, len(member))
	result := make([]string, 0, len(member))

	for _, m := range member {
		if _, exists := unique[m.MemberHash]; exists {
			continue
		}
		unique[m.MemberHash] = struct{}{}
		result = append(result, m.MemberHash)
	}

	r.chatRoomMemberMap[roomKey] = result

	log.Printf(
		"[PutChatRoomMember] regist. roomKey=%s members=%v",
		roomKey,
		r.chatRoomMemberMap[roomKey],
	)

}
