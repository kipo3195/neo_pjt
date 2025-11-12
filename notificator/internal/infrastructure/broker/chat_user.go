package broker

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ChatUser struct {
	mu   sync.RWMutex
	ch   chan BrokerMessage
	conn *websocket.Conn // 각 유저는 1개의 커넥션만 유지
}
