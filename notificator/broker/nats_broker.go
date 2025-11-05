package broker

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	Conn      *nats.Conn
	mu        sync.RWMutex         // 명시적으로 초기화하지 않아도 자동으로 초기화됩니다.
	ChatRooms map[string]*ChatRoom // 일반 채팅용 채널 관리용 map
}

// 추후 entity로 변경
type SimpleMessage struct {
	roomId string
	data   []byte
}

func (m *SimpleMessage) Data() []byte {
	return m.data
}

func (mb *NatsBroker) SubscribeChatRoom(roomId string, userId string, conn *websocket.Conn) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// 맵에 해당 방이 있는지 체크, 없으면 메모리에 추가를 위한 정보 생성
	room, exists := mb.ChatRooms[roomId]
	if !exists {
		room = &ChatRoom{
			ch:      make(chan BrokerMessage, 100),
			clients: make(map[string]*websocket.Conn),
		}

		mb.ChatRooms[roomId] = room

		// fan-out goroutine 시작
		go func(r *ChatRoom) {
			for msg := range r.ch {
				r.mu.RLock()
				for _, c := range r.clients {
					go func(c *websocket.Conn) {
						c.WriteMessage(websocket.TextMessage, msg.Data())
					}(c)
				}
				r.mu.RUnlock()
			}
		}(room)
	}

	// 구독 사용자 추가
	room.mu.Lock()
	room.clients[userId] = conn
	room.mu.Unlock()
}

func (mb *NatsBroker) UnSubscribeChatRoom(roomId string, userId string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if room, exists := mb.ChatRooms[roomId]; exists {

		// 유저 삭제
		room.mu.Lock()
		delete(room.clients, userId)
		roomEmpty := len(room.clients) == 0
		room.mu.Unlock()

		// 현재 방을 구독하는 유저가 없을때만
		if roomEmpty {
			close(room.ch)               // 채널 닫기
			delete(mb.ChatRooms, roomId) // 방의 채널을 맵에서 제거
		}
	}
}

func (mb *NatsBroker) PublishToChatRoom(roomId string, data []byte) error {
	mb.mu.RLock()
	room, exists := mb.ChatRooms[roomId]
	defer mb.mu.RUnlock()

	if !exists {
		fmt.Printf("room %s is not exist", roomId)
		return errors.New("room is not exist")
	}

	select {
	case room.ch <- &SimpleMessage{roomId: roomId, data: data}:
		fmt.Printf("sending success room %s ", string(roomId))
		return nil
	default:
		fmt.Printf("room %s channel is full", roomId)
		return errors.New("room channel is full") // fan-out 병목 가능성 알림
	}
}
