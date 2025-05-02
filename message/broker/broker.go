package broker

type Broker interface {
	Subscribe(roomId string) (Subscription, chan BrokerMessage, error)
	Publish(roomId string, data []byte) error
}

type Subscription interface {
	Unsubscribe() error
}

type BrokerMessage interface {
	Data() []byte
}
