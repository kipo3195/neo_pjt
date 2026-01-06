package storage

import (
	"log"
	"sync"
)

type chatRoomStorage struct {
	mu sync.RWMutex
	//chatUserConnectMap map[string]*websocket.Conn
	// chatRoomMemberMap  map[string][]string -> 이 구조는 방에 참여자가 많아질수록 방의 수 * 참여자의 수만큼 반복해야하므로.. 개선함
	chatRoomMemberMap map[string]map[string]struct{} // roomKey : 참여자SET 의 형태를 취함.  -> 채팅방 수신시 사용자에게 write 하기 위한 용도
	memberChatRoomMap map[string]map[string]struct{} // 참여자 : roomkey SET 의 형태를 취함. -> 소켓 disconnect시 내가 참여 중인 방을 정리하기 위한용도
}

type ChatRoomStorage interface {
	// GetChatConnect(userHash string) *websocket.Conn
	// RemoveChatConnect(userHash string)
	//PutChatConnect(userHash string, conn *websocket.Conn, c chan []byte)
	GetChatRoomMember(roomKey string) []string
	PutChatRoomMember(roomKey string, member []string)
	InitMyRoom(roomKey []string, userHash string)
	CleanUpMyRoom(userHash string)
}

func NewChatRoomStorage() ChatRoomStorage {
	return &chatRoomStorage{
		//chatUserConnectMap: make(map[string]*websocket.Conn),
		chatRoomMemberMap: make(map[string]map[string]struct{}),
		memberChatRoomMap: make(map[string]map[string]struct{}),
	}
}

// func (r *chatUserStorage) GetChatConnect(userHash string) *websocket.Conn {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()

// 	return r.chatUserConnectMap[userHash]
// }
// func (r *chatUserStorage) RemoveChatConnect(userHash string) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	connect := r.chatUserConnectMap[userHash]

//		if connect != nil {
//			delete(r.chatUserConnectMap, userHash)
//		}
//		log.Println("[RemoveChatConnect] userHash : ", userHash)
//	}
// func (r *chatUserStorage) PutChatConnect(userHash string, conn *websocket.Conn, c chan []byte) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	r.chatUserConnectMap[userHash] = conn
// 	log.Println("[PutChatConnect] userHash : ", userHash)

// }

func (r *chatRoomStorage) GetChatRoomMember(roomKey string) []string {
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

func (r *chatRoomStorage) PutChatRoomMember(roomKey string, member []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 해당 방의 멤버 맵 초기화
	newMembers := make(map[string]struct{}, len(member))
	for _, m := range member {
		newMembers[m] = struct{}{}
		r.memberChatRoomMap[m][roomKey] = struct{}{}
	}

	r.chatRoomMemberMap[roomKey] = newMembers

	log.Printf(
		"[PutChatRoomMember] roomKey=%s members=%v",
		roomKey,
		r.chatRoomMemberMap[roomKey],
	)

}

func (r *chatRoomStorage) InitMyRoom(roomKey []string, userHash string) {

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

		// 참여자가 없으면 생성
		if _, exists := r.memberChatRoomMap[userHash]; !exists {
			r.memberChatRoomMap[userHash] = make(map[string]struct{})
		}

		// Set 구조이므로 중복 체크를 위해 루프를 돌 필요가 없음 (O(1))
		r.chatRoomMemberMap[rk][userHash] = struct{}{}
		r.memberChatRoomMap[userHash][rk] = struct{}{}

		log.Printf(
			"[InitMyRoom] roomKey=%s members=%v",
			rk,
			r.chatRoomMemberMap[rk],
		)
	}
}

func (r *chatRoomStorage) CleanUpMyRoom(userHash string) {

	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. 해당 유저가 속한 모든 roomKey 목록을 가져옴 (memberChatRoomMap 이용)
	rooms, exists := r.memberChatRoomMap[userHash]
	if !exists {
		log.Printf("[CleanUpMyRoom] No active rooms found for user: %s", userHash)
		return
	}

	// 2. 찾아낸 각 roomKey를 순회하며 chatRoomMemberMap에서 유저를 제거
	for roomKey := range rooms {
		if members, ok := r.chatRoomMemberMap[roomKey]; ok {
			delete(members, userHash) // 방 멤버 목록에서 나를 삭제

			// (선택 사항) 만약 방에 아무도 남지 않았다면 방 자체를 삭제하여 메모리 절약
			if len(members) == 0 {
				delete(r.chatRoomMemberMap, roomKey)
				log.Printf("[CleanUpMyRoom] Room %s is now empty and removed", roomKey)
			}
		}
	}

	// 3. 마지막으로 유저 기준 맵에서도 해당 유저 데이터를 완전히 삭제
	delete(r.memberChatRoomMap, userHash)

	log.Printf("[CleanUpMyRoom] complete for user: %s", userHash)

}
