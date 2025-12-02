package entity

import "github.com/gorilla/websocket"

type SendConnectionEntity struct {
	UserHash string
	Conn     *websocket.Conn
	Chan     chan interface{}
}

func MakeSendConnectionEntity(userHash string, conn *websocket.Conn, c chan interface{}) SendConnectionEntity {

	return SendConnectionEntity{
		UserHash: userHash,
		Conn:     conn,
		Chan:     c,
	}

}
