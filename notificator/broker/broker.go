package broker

import "github.com/gorilla/websocket"

type Broker interface {
	SubscribeChatRoom(roomId string, userId string, conn *websocket.Conn)
	UnSubscribeChatRoom(roomId string, userId string)
	PublishToChatRoom(roomId string, data []byte) error
}

type Subscription interface {
	Unsubscribe() error
}

type BrokerMessage interface {
	Data() []byte
}
