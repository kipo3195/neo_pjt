package broker

import (
	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	Connector *nats.Conn
	//mu        sync.RWMutex         // 명시적으로 초기화하지 않아도 자동으로 초기화됩니다.
	ChatUsers map[string]*ChatUser // 일반 채팅용 채널 관리용 map
}

// 추후 entity로 변경
// type SimpleMessage struct {
// 	roomId string
// 	data   []byte
// }

// func (m *SimpleMessage) Data() []byte {
// 	return m.data
// }

// func (mb *NatsBroker) PublishToChatUser(userHash string, contents string) error {
// 	subject := fmt.Sprintf("chat.%s", userHash)

// 	data, err := json.Marshal(contents)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal broker message: %w", err)
// 	}

// 	if err := mb.Nc.Publish(subject, data); err != nil {
// 		return fmt.Errorf("failed to publish NATS message: %w", err)
// 	}

// 	return nil
// }
