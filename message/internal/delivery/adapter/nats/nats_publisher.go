package nats

import (
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{nc: nc}
}

func (p *Publisher) Publish(subject string, data []byte) error {
	return p.nc.Publish(subject, data)
}
