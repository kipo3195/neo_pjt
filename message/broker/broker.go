package broker

type Broker interface {
	SubscribeChatRoom(roomId string) (chan BrokerMessage, error)
	UnsubscribeChatRoom(roomId string)
	PublishToChatRoom(roomId string, data []byte) error
}

type Subscription interface {
	Unsubscribe() error
}

type BrokerMessage interface {
	Data() []byte
}
