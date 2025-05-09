package broker

import (
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	Conn          *nats.Conn
	mu            sync.RWMutex
	Subscriptions map[string]chan BrokerMessage // roomKey : channel
}

type SimpleMessage struct {
	roomId string
	data   []byte
}

func (m *SimpleMessage) Data() []byte {
	return m.data
}
func (mb *NatsBroker) Subscribe(room string) (chan BrokerMessage, error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// subscriptions 맵에 해당 방이 있는지 체크
	ch, exists := mb.Subscriptions[room]
	if !exists {
		ch = make(chan BrokerMessage, 100)
		mb.Subscriptions[room] = ch
	}
	return ch, nil
}

func (mb *NatsBroker) Unsubscribe(room string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if ch, exists := mb.Subscriptions[room]; exists {
		close(ch)                      // 채널 닫기
		delete(mb.Subscriptions, room) // 방의 채널을 맵에서 제거
	}
}

// NATS용 채널과 구독해제용 Subscription 리턴
// func (b *NatsBroker) Subscribe(roomId string) (Subscription, chan BrokerMessage, error) {
// 	natsMsgChan := make(chan *nats.Msg, 64)       // ChanSubscribe용 NATS 채널
// 	brokerMsgChan := make(chan BrokerMessage, 64) // 우리가 핸들링할 BrokerMessage 채널

// 	sub, err := b.Conn.ChanSubscribe(roomId, natsMsgChan)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// 변환용 고루틴
// 	// 이 방법은 ChanSubscribe 그대로 활용하면서 BrokerMessage 타입으로 채널 추상화.
// 	go func() {
// 		for msg := range natsMsgChan {
// 			brokerMsgChan <- &NatsMessage{msg: msg}
// 		}
// 	}()

// 	return &NatsSubscription{sub: sub}, brokerMsgChan, nil
// }

type NatsMessage struct {
	msg *nats.Msg
}

func (n *NatsMessage) Data() []byte {
	return n.msg.Data
}

type NatsSubscription struct {
	sub *nats.Subscription
}

func (n *NatsSubscription) Unsubscribe() error {
	return n.sub.Unsubscribe()
}

func (mb *NatsBroker) Publish(roomId string, data []byte) error {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	msg := &SimpleMessage{
		roomId: roomId,
		data:   data,
	}

	if ch, exists := mb.Subscriptions[roomId]; exists {
		ch <- msg
	}
	return nil
}
