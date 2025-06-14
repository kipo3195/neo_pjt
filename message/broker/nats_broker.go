package broker

import (
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	Conn                  *nats.Conn
	mu                    sync.RWMutex
	ChatRoomSubscriptions map[string]chan BrokerMessage // 일반 채팅용 채널 관리용 map
}

// 추후 entity로 변경
type SimpleMessage struct {
	roomId string
	data   []byte
}

func (m *SimpleMessage) Data() []byte {
	return m.data
}

func (mb *NatsBroker) SubscribeChatRoom(room string) (chan BrokerMessage, error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// 맵에 해당 방이 있는지 체크, 구독하는 사용자는 모두 같은 채널을 받음
	ch, exists := mb.ChatRoomSubscriptions[room]
	if !exists {
		ch = make(chan BrokerMessage, 100)
		mb.ChatRoomSubscriptions[room] = ch
	}
	return ch, nil
}

func (mb *NatsBroker) UnsubscribeChatRoom(room string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if ch, exists := mb.ChatRoomSubscriptions[room]; exists {
		close(ch)                              // 채널 닫기
		delete(mb.ChatRoomSubscriptions, room) // 방의 채널을 맵에서 제거
	}
}

func (mb *NatsBroker) PublishToChatRoom(roomId string, data []byte) error {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	msg := &SimpleMessage{
		roomId: roomId,
		data:   data,
	}

	if ch, exists := mb.ChatRoomSubscriptions[roomId]; exists {
		ch <- msg
	}
	return nil
}

// type NatsMessage struct {
// 	msg *nats.Msg
// }

// func (n *NatsMessage) Data() []byte {
// 	return n.msg.Data
// }
