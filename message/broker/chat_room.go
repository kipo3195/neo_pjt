package broker

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	ch      chan BrokerMessage         // 이 방에서 발생하는 메시지를 받아올 단일 채널
	clients map[string]*websocket.Conn // userId → WebSocket 연결
	mu      sync.RWMutex               // clients 맵 보호용 락
}
