package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"notificator/internal/application/orchestrator"
	"notificator/internal/application/usecase/input"
	"notificator/internal/consts"
	"notificator/internal/delivery/adapter"
	"notificator/internal/delivery/dto/notificatorService"
	commonConsts "notificator/pkg/consts"
	"notificator/pkg/response"

	"github.com/gorilla/websocket"
)

type NotificatorServiceHandler struct {
	svc *orchestrator.NotificatorService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewNotificatorServiceHandler(svc *orchestrator.NotificatorService) *NotificatorServiceHandler {
	return &NotificatorServiceHandler{
		svc: svc,
	}
}

func (h *NotificatorServiceHandler) NotificatorConnect(w http.ResponseWriter, r *http.Request) {

	/* 웹소켓 연결 수립 */
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		/* 웹소켓 업그레이드 에러 */
		log.Println("Notificator service WebSocket upgrade error:", err)
		response.SendErrorWs(conn, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 함수 종료시 소켓 연결 종료
	defer conn.Close()

	/* header에 있는 AT를 파싱하여 사용자 정보 체크 */
	userHash := r.Context().Value(consts.USER_HASH)
	log.Println("Notificator service connect success! userHash :", userHash)

	/* 연결 성공에 대한 response 처리 */
	res := notificatorService.NotificatorConnectResponse{
		UserHash: userHash.(string),
	}
	response.SendSuccess(conn, res)

	// 쓰기 (server -> client) 채널 생성 후 메모리 저장, 쓰기 고루틴 시작
	go h.svc.SocketSender.SaveConnection(conn, userHash.(string))

	// 읽기 (client -> server)
	for {
		// 메시지는 반복해서 수신, ReadMessage는 블로킹 함수
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Notificator service Read msg error:", err)
			break
		}

		// type 파싱
		// type 에 따라 data 부분을 다르게 언마샬 하려면, 바깥은 공통 envelope, 안쪽은 json.RawMessage 로 둡니다.
		var req notificatorService.NotificatorConnectRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("Notificator service websocket message error:", err)
			return
		}

		// 여기서 각각의 usecase를 활용한 처리

		switch req.Type {

		case consts.AUTH:

		case consts.LOGIN:
			var body notificatorService.LoginDto
			if err := json.Unmarshal(req.Data, &body); err != nil {
				log.Println("[LOGIN] data unmarshal error:", err)
				return
			}

			input := adapter.MakeLoginInput(body.Uuid, body.DeviceType)
			h.svc.Login.LoginProcess(input)

		// 채팅, 쪽지 구독이 필요한가?
		case consts.CHAT:
			var input input.ChatConnectInput
			if err := json.Unmarshal(msg, &input); err == nil {
				//h.svc.Chat.SubscribeChat(input, conn)
				log.Println("Notificator service chat subscribe success.")
			}

		case consts.NOTE:
			var input input.NoteConnectInput
			if err := json.Unmarshal(msg, &input); err == nil {
				h.svc.Note.SubscribeNote(input, conn)
				log.Println("Notificator service note subscribe success.")
			}

		default:
			log.Println("unknown message type:", req.Type)
			return
		}

		// 연결 종료 로직 추가 필요

	}

	log.Println("Notificator service websocket close. ")
}
