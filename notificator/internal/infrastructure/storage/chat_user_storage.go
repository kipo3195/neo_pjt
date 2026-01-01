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
	// chatRoomMemberMap  map[string][]string -> 이 구조는 방에 참여자가 많아질수록 방의 수 * 참여자의 수만큼 반복해야하므로.. 개선함
	chatRoomMemberMap map[string]map[string]struct{} // roomKey : 참여자SET 의 형태를 취함.
}

type ChatUserStorage interface {
	GetChatConnect(userHash string) *websocket.Conn
	RemoveChatConnect(userHash string)
	PutChatConnect(userHash string, conn *websocket.Conn, c chan []byte)
	GetChatRoomMember(roomKey string) []string
	PutChatRoomMember(roomKey string, member []entity.ChatRoomMemberEntity)
	InitMyRoom(roomKey []string, userHash string)
}

func NewChatUserStorage() ChatUserStorage {
	return &chatUserStorage{
		chatUserConnectMap: make(map[string]*websocket.Conn),
		chatRoomMemberMap:  make(map[string]map[string]struct{}),
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
	log.Println("[GetChatRoomMember] ", roomKey)

	members, exists := r.chatRoomMemberMap[roomKey]
	if !exists || len(members) == 0 {
		// 채팅방이 생성되지 않았거나, 채팅방에 참여중인 member가 소켓 연결 되지 않은 상태.
		log.Println("[GetChatRoomMember] roomKey is not regist or member is not init.")
		return nil
	}

	// Map을 Slice로 변환하여 반환
	result := make([]string, 0, len(members))
	for memberHash := range members {
		result = append(result, memberHash)
	}

	return result
}

func (r *chatUserStorage) PutChatRoomMember(roomKey string, member []entity.ChatRoomMemberEntity) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 해당 방의 멤버 맵 초기화
	newMembers := make(map[string]struct{}, len(member))
	for _, m := range member {
		newMembers[m.MemberHash] = struct{}{}
	}

	r.chatRoomMemberMap[roomKey] = newMembers

	log.Printf(
		"[PutChatRoomMember] roomKey=%s members=%v",
		roomKey,
		r.chatRoomMemberMap[roomKey],
	)

}

func (r *chatUserStorage) InitMyRoom(roomKey []string, userHash string) {

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, rk := range roomKey {
		if rk == "" {
			continue
		}

		// 방이 없으면 생성
		if _, exists := r.chatRoomMemberMap[rk]; !exists {
			r.chatRoomMemberMap[rk] = make(map[string]struct{})
		}

		// Set 구조이므로 중복 체크를 위해 루프를 돌 필요가 없음 (O(1))
		r.chatRoomMemberMap[rk][userHash] = struct{}{}

		log.Printf(
			"[InitMyRoom] roomKey=%s members=%v",
			rk,
			r.chatRoomMemberMap[rk],
		)
	}
}
