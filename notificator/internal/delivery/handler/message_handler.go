package handler

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"notificator/internal/consts"
// 	"notificator/internal/infrastructure/broker"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

type MessageHandler struct {
	// cu usecase.ChatUsecase
	// nu usecase.NoteUsecase
}

// // 웹소켓은 하나의 핸들러에서 처리 단, useCase를 여러개 둘 수 있음.
// func NewMessageHandler(cu usecase.ChatUsecase, au usecase.AuthUsecase, nu usecase.NoteUsecase, mb broker.Broker) *MessageHandler {
// 	return &MessageHandler{
// 		cu: cu,
// 		nu: nu,
// 	}
// }

// func parseToResponseByte(response *dto.Response) []byte {

// 	data, err := json.Marshal(response)
// 	if err != nil {
// 		return []byte(`{"error":"server check."}`)
// 	}
// 	return data
// }

// // MessageHandler의 메소드 SetupRoutes에서 등록됨.
// func (h *MessageHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

// 	// response
// 	res := &dto.Response{}

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("WebSocket upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// 인증 처리 flag
// 	authenticated := false

// 	for {
// 		// 메시지는 반복해서 수신, ReadMessage는 블로킹 함수
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read msg error:", err)
// 		}

// 		/* 페이로드 파싱 예는 entity여야 하나? cmd, data 형식으로, cmd에 따라 data를 어떤 entity로 파싱할지 처리?*/
// 		var data map[string]interface{}
// 		json.Unmarshal(msg, &data)
// 		log.Println("웹소켓 메시지 수신시 최초 로깅 : ", data)
// 		msgType, ok := data[consts.TYPE].(string)
// 		if !ok {
// 			res.Result = consts.ERROR
// 			res.Data = &dto.ErrorResponse{
// 				Code:    consts.E_101,
// 				Message: consts.E_101_MSG,
// 			}
// 			conn.WriteMessage(websocket.TextMessage, parseToResponseByte(res))
// 			continue
// 		}

// 		/*최초 인증 처리*/
// 		if !authenticated {
// 			switch msgType {
// 			case consts.AUTH:
// 				result, err := h.au.HandleAuth(conn, data)
// 				if !result && err != nil {
// 					res.Result = consts.ERROR
// 					res.Data = &dto.ErrorResponse{
// 						Code:    consts.MESSAGE_F001,
// 						Message: consts.MESSAGE_F001_MSG,
// 					}
// 					conn.WriteMessage(websocket.TextMessage, parseToResponseByte(res))
// 					return
// 				} else {
// 					// 성공시에는 성공했다고 필요할까? - 클라이언트 논의
// 					// conn.WriteMessage(websocket.TextMessage, []byte(`{"success":"auth_ok"}`))
// 					log.Println("연결 성공 !")
// 					authenticated = true
// 					continue
// 				}
// 			default:
// 				// 이걸 받는다면 인증되지 않았으므로 소켓 연결을 다시 해야 될 듯.
// 				res.Result = consts.ERROR
// 				res.Data = &dto.ErrorResponse{
// 					Code:    consts.E_106,
// 					Message: consts.E_106_MSG,
// 				}
// 				conn.WriteMessage(websocket.TextMessage, parseToResponseByte(res))
// 				return
// 			}
// 		} else {

// 			// 인증 이후에는 웹소켓 수신 메시지를 가지고 이쪽으로
// 			switch msgType {
// 			case consts.CHAT:
// 				h.cu.HandleChat(conn, data)
// 			case consts.NOTE:
// 				h.nu.HandleNote(conn, data)
// 			}
// 		}

// 	}

// }
