package broker

import (
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	Connector *nats.Conn
	mu        sync.RWMutex // 명시적으로 초기화하지 않아도 자동으로 초기화됩니다.
}

func (mb *NatsBroker) GetConnector() interface{} {
	return NatsBroker{}
}

// func (mb *NatsBroker) SubscribeChatRoom(roomId string, userId string, conn *websocket.Conn) {
// 	mb.mu.Lock()
// 	defer mb.mu.Unlock()

// 	// 맵에 해당 방이 있는지 체크, 없으면 메모리에 추가를 위한 정보 생성
// 	room, exists := mb.ChatRooms[roomId]
// 	if !exists {
// 		room = &ChatRoom{
// 			ch:      make(chan BrokerMessage, 100),
// 			clients: make(map[string]*websocket.Conn),
// 		}

// 		mb.ChatRooms[roomId] = room

// 		// fan-out goroutine 시작
// 		go func(r *ChatRoom) {
// 			for msg := range r.ch {
// 				r.mu.RLock()
// 				for _, c := range r.clients {
// 					go func(c *websocket.Conn) {
// 						c.WriteMessage(websocket.TextMessage, msg.Data())
// 					}(c)
// 				}
// 				r.mu.RUnlock()
// 			}
// 		}(room)
// 	}

// 	// 구독 사용자 추가
// 	room.mu.Lock()
// 	room.clients[userId] = conn
// 	room.mu.Unlock()
// }

// func (mb *NatsBroker) UnSubscribeChatRoom(roomId string, userId string) {
// 	mb.mu.Lock()
// 	defer mb.mu.Unlock()

// 	if room, exists := mb.ChatRooms[roomId]; exists {

// 		// 유저 삭제
// 		room.mu.Lock()
// 		delete(room.clients, userId)
// 		roomEmpty := len(room.clients) == 0
// 		room.mu.Unlock()

// 		// 현재 방을 구독하는 유저가 없을때만
// 		if roomEmpty {
// 			close(room.ch)               // 채널 닫기
// 			delete(mb.ChatRooms, roomId) // 방의 채널을 맵에서 제거
// 		}
// 	}
// }

// func (mb *NatsBroker) PublishToChatRoom(roomId string, data []byte) error {
// 	mb.mu.RLock()
// 	room, exists := mb.ChatRooms[roomId]
// 	defer mb.mu.RUnlock()

// 	if !exists {
// 		fmt.Printf("room %s is not exist", roomId)
// 		return errors.New("room is not exist")
// 	}

// 	select {
// 	case room.ch <- &SimpleMessage{roomId: roomId, data: data}:
// 		fmt.Printf("sending success room %s ", string(roomId))
// 		return nil
// 	default:
// 		fmt.Printf("room %s channel is full", roomId)
// 		return errors.New("room channel is full") // fan-out 병목 가능성 알림
// 	}
// }

// func (mb *NatsBroker) SubscribeChat(userHash string, conn *websocket.Conn) {
// 	mb.mu.Lock()
// 	defer mb.mu.Unlock()

// 	// 기존 구독 존재 확인
// 	user, exists := mb.ChatUsers[userHash]
// 	if !exists {
// 		user = &ChatUser{
// 			ch:   make(chan BrokerMessage, 100),
// 			conn: conn,
// 		}
// 		mb.ChatUsers[userHash] = user

// 		// fan-out goroutine (1개 유저만 상대하므로 단순)
// 		go func(r *ChatUser, hash string) {
// 			for msg := range r.ch {
// 				if r.conn != nil {
// 					_ = r.conn.WriteMessage(websocket.TextMessage, msg.Data())
// 				}
// 			}
// 		}(user, userHash)

// 		// NATS 구독 (userHash 전용)
// 		subject := fmt.Sprintf("chat.%s", userHash)
// 		_, err := mb.Connector.Subscribe(subject, func(m *nats.Msg) {
// 			var bm BrokerMessage
// 			if err := json.Unmarshal(m.Data, &bm); err == nil {
// 				select {
// 				case user.ch <- bm:
// 				default:
// 					log.Printf("[%s] channel full, dropping message", user)
// 				}
// 			}
// 		})
// 		if err != nil {
// 			log.Println("NATS subscribe error:", err)
// 		}

// 		log.Printf("User %s subscribed to %s", userHash, subject)
// 	} else {
// 		// 이미 존재하면 conn 갱신 (재접속 등)
// 		user.conn = conn
// 	}
// }
