package broker

type Broker interface {
	Subscribe(roomId string) (chan BrokerMessage, error)
	Unsubscribe(roomId string)
	Publish(roomId string, data []byte) error
}

type Subscription interface {
	Unsubscribe() error
}

type BrokerMessage interface {
	Data() []byte
}
