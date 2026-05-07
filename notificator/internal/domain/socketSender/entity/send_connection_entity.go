package entity

import "github.com/gorilla/websocket"

type SendConnectionEntity struct {
	Conn *websocket.Conn // WebSocket Conn은 복사되면 안 되는 타입
	Chan chan interface{}
}

func MakeSendConnectionEntity(userHash string, conn *websocket.Conn, c chan interface{}) *SendConnectionEntity {

	return &SendConnectionEntity{
		Conn: conn,
		Chan: c,
	}

}
