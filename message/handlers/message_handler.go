package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"message/broker"
	"message/usecases"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

const (
	TYPE = "type" // 이벤트의 대분류
	CMD  = "cmd"  // 이벤트 type의 세부 타입
	AUTH = "auth"
	CHAT = "chat"
	NOTE = "note"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageHandler struct {
	au usecases.AuthUsecase
	cu usecases.ChatUsecase
	mb broker.Broker
}

// 웹소켓은 하나의 핸들러에서 처리 단, useCase를 여러개 둘 수 있음.
func NewMessageHandler(cu usecases.ChatUsecase, au usecases.AuthUsecase, mb broker.Broker) *MessageHandler {
	return &MessageHandler{
		cu: cu,
		au: au,
		mb: mb,
	}
}

// 덕타이핑스러운 처리
// 구독 정보 interface Unsubscribe를 가짐
//
//	nats의 NatsSubscription은 Unsubscribe를 가지고 있기 때문에 처리가능함.
// type Subscription interface {
// 	Unsubscribe() error
// }

// /*NATS용 래퍼 정의*/
// type NatsSubscription struct {
// 	sub *nats.Subscription
// }

// /*NATS용 구독해제 메소드 정의 */
// func (n *NatsSubscription) Unsubscribe() error {
// 	return n.sub.Unsubscribe()
// }

type NatsMessage struct {
	msg *nats.Msg
}

func (n *NatsMessage) Data() []byte {
	return n.msg.Data
}

/*NATS용 브로커 (덕타이핑)*/
// type NatsBroker struct {
// 	conn *nats.Conn
// }

/*NATS용 구독 메소드 정의 */
// func (b *NatsBroker) Subscribe(roomId string) (Subscription, error) {
// 	msgChan := make(chan *nats.Msg, 64)
// 	natsSub, err := b.conn.ChanSubscribe(roomId, msgChan)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &NatsSubscription{sub: natsSub}, nil
// }

type incomming struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

// MessageHandler의 메소드 SetupRoutes에서 등록됨.
func (h *MessageHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// 인증 처리 flag
	authenticated := false

	var sub broker.Subscription // 구독 정보

	for {
		// 메시지는 반복해서 수신
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read msg error:", err)
		}

		/* 페이로드 검증 */
		var data map[string]interface{}
		json.Unmarshal(msg, &data)
		fmt.Println("웹소켓 메시지 수신시 최초 로깅 : ", data)
		msgType := data[TYPE].(string)

		// 최초 인증 처리 -> 모듈화? 방안 모색 TODO
		if !authenticated {
			if msgType == AUTH {
				token := data["token"].(string)

				flag, err := h.au.AuthenticateToken(token)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"auth_failed"}`))
					// defer conn.Close() 호출 되므로 연결 종료.
					break
				}
				authenticated = flag
				conn.WriteMessage(websocket.TextMessage, []byte(`{"success":"auth_ok"}`))
				// continue하지않으면 아래로 내려가 버리기 때문에 인증완료 후에는 다시 메시지 수신(블로킹)
				continue
			} else {
				// 인증 안 된 상태에서는 다른 메시지 못 보내게 함
				// defer conn.Close() 호출 되므로 연결 종료.
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"unauthorized"}`))
				break
			}

		}

		fmt.Println("인증 완료된 계정 :", data["userId"])
		userId := data["userId"]
		// 채팅방을 구독하더라도 쪽지를 보낼 수 있어야 한다.
		cmd := data[CMD].(string)
		payload := data["payload"].(map[string]interface{})
		if msgType == CHAT {

			//msgChan := make(chan *nats.Msg, 64) // 채널은 생성, 채널에 데이터를 구독하는건 joinRoom에서
			var msgChan chan broker.BrokerMessage

			roomId := payload["roomId"].(string)
			if cmd == "joinRoom" {
				// 채널 생성

				// interface였던 브로커를 꺼내기 위함. -> 다른 mb라면???
				mbSub, bm, err := h.mb.Subscribe(roomId)
				if err != nil {
					log.Println("NATS subscribe error:", err)
					return
				}
				// 래퍼로 감싸서 인터페이스로 저장
				sub = mbSub
				msgChan = bm
				fmt.Printf("roomId : %s 구독 하였습니다.\n", roomId)

				go func() {
					for {
						m := <-msgChan
						var initReq incomming
						_ = json.Unmarshal(m.Data(), &initReq)
						// fmt.Println("채널에서 수신받은 데이터 :", initReq.Content)
						if initReq.Sender == userId {
							continue
						}
						conn.WriteMessage(websocket.TextMessage, m.Data())
					}
				}()
			} else if cmd == "joinRoomCancel" {
				// 구독을 해제 하였을때도 sendMessage가 동작해야하나? TODO
				if sub != nil {
					sub.Unsubscribe() // 인터페이스 메서드 호출 (OK!)
					sub = nil
					fmt.Printf("roomId : %s 구독해제 하였습니다.\n", roomId)
				}
			} else if sub != nil && cmd == "sendMessage" {
				// 별도의 고루틴으로 처리?
				// 채팅방을 구독하지 않는 상태에서는 메시지를 보낼 수 없는 것이 맞지않을까?
				content := payload["content"].(string)
				payload := map[string]string{"content": content, "sender": userId.(string)}
				jsonBytes, _ := json.Marshal(payload)

				if err := h.mb.Publish(roomId, jsonBytes); err != nil {
					log.Println("Publish error:", err)
				}
			}

		} else if msgType == NOTE {
			if cmd == "sendMessage" {

			}

		}

	}

}
