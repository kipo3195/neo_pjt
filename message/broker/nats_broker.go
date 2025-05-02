package broker

import "github.com/nats-io/nats.go"

type NatsBroker struct {
	Conn *nats.Conn
}

// NATS용 채널과 구독해제용 Subscription 리턴
func (b *NatsBroker) Subscribe(roomId string) (Subscription, chan BrokerMessage, error) {
	natsMsgChan := make(chan *nats.Msg, 64)       // ChanSubscribe용 NATS 채널
	brokerMsgChan := make(chan BrokerMessage, 64) // 우리가 핸들링할 BrokerMessage 채널

	sub, err := b.Conn.ChanSubscribe(roomId, natsMsgChan)
	if err != nil {
		return nil, nil, err
	}

	// 변환용 고루틴
	// 이 방법은 ChanSubscribe 그대로 활용하면서 BrokerMessage 타입으로 채널 추상화.
	go func() {
		for msg := range natsMsgChan {
			brokerMsgChan <- &NatsMessage{msg: msg}
		}
	}()

	return &NatsSubscription{sub: sub}, brokerMsgChan, nil
}

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

func (b *NatsBroker) Publish(roomId string, data []byte) error {
	return b.Conn.Publish(roomId, data)
}
