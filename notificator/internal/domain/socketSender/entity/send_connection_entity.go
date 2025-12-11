package entity

import "github.com/gorilla/websocket"

type SendConnectionEntity struct {
	UserHash string
	Conn     *websocket.Conn // WebSocket Conn은 복사되면 안 되는 타입
	Chan     chan interface{}
}

func MakeSendConnectionEntity(userHash string, conn *websocket.Conn, c chan interface{}) *SendConnectionEntity {

	return &SendConnectionEntity{
		UserHash: userHash,
		Conn:     conn,
		Chan:     c,
	}

}
