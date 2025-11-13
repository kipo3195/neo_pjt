package broker

type Broker interface {
	// SubscribeChatRoom(roomId string, userId string, conn *websocket.Conn)
	// UnSubscribeChatRoom(roomId string, userId string)
	// PublishToChatRoom(roomId string, data []byte) error
	// SubscribeChat(userhash string, conn *websocket.Conn)
	GetConnector() interface{}
}

type Subscription interface {
	Unsubscribe() error
}

type BrokerMessage interface {
	Data() []byte
}
