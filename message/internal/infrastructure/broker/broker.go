package broker

type Broker interface {
	//PublishToChatUser(userHash string, data string) error
}

type Subscription interface {
	Unsubscribe() error
}

type BrokerMessage interface {
	Data() []byte
}
