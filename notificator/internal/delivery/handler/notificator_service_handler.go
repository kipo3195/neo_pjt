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
	"notificator/internal/infrastructure/config"
	commonConsts "notificator/pkg/consts"
	"notificator/pkg/response"
	"time"

	"github.com/gorilla/websocket"
)

type NotificatorServiceHandler struct {
	svc             *orchestrator.NotificatorService
	websocketConfig config.WebsocketConnectionConfig
}

func NewNotificatorServiceHandler(svc *orchestrator.NotificatorService, websocketConfig config.WebsocketConnectionConfig) *NotificatorServiceHandler {
	return &NotificatorServiceHandler{
		svc:             svc,
		websocketConfig: websocketConfig,
	}
}

// Upgrader 설정과 보안 (CheckOrigin)
// 보안을 위해 교차 출처(CORS)를 차단
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *NotificatorServiceHandler) NotificatorConnect(w http.ResponseWriter, r *http.Request) {

	var pongWait = time.Duration(h.websocketConfig.PongWait) * time.Second // 클라이언트의 응답을 기다릴 최대 시간

	/* 웹소켓 연결 수립 */
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		/* 웹소켓 업그레이드 에러 */
		log.Println("Notificator service WebSocket upgrade error:", err)
		response.SendErrorWs(conn, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	/* 함수 종료시 소켓 연결 종료 */
	defer conn.Close()

	/* header에 있는 AT를 파싱하여 사용자 정보 체크 */
	user := r.Context().Value(consts.USER_HASH)
	log.Println("Notificator service connect request userHash :", user)

	if user == nil || user == "" {
		log.Println("Notificator service connect error. userHash invalid.")

		// todo error msg
		return
	}
	userHash := user.(string)

	/* 쓰기 (server -> client) 채널 생성 후 메모리 저장, 쓰기 고루틴 시작 */
	go h.svc.SocketSender.SaveConnection(conn, userHash, h.websocketConfig)

	/* 내가 참여중인 방 notificator 서비스 메모리에 로딩 처리 시작 */
	err = h.svc.Chat.SubscribeChat(userHash)
	if err != nil {
		log.Println("Notificator service connect error. Subscribe chat room error.")

		// todo error msg
		return
	}

	/* pong 처리 */
	// pong을 for 문안에서 처리하지않아도 되는 이유
	// pong을 ReadMessage()가 반환하는 데이터가 아니기 때문에 switch문 안에서 처리하실 필요가 없습니다.
	// ReadMessage()는 텍스트 메시지, 바이너리 메시지 등을 처리하지만 pong 메시지는 별도의 핸들러에서 처리됩니다.
	// 따라서 pong 메시지를 받으면 설정된 PongHandler가 자동으로 호출되어 처리됩니다.

	// ReadDeadline 설정 (지금부터 설정한 시간 안에 아무 메시지나 Pong이 와야 함)
	conn.SetReadDeadline(time.Now().Add(pongWait))

	// PongHandler 설정: 클라이언트로부터 Pong 메시지를 받으면 호출됨
	conn.SetPongHandler(func(string) error {
		// Pong을 받았으므로 마감 기한을 다시 뒤로 연장
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	/* 연결 성공에 대한 response 처리 */
	res := notificatorService.NotificatorConnectResponse{
		UserHash: userHash,
	}
	response.SendSuccess(conn, res)

	// 읽기 (client -> server)
	for {
		// 메시지는 반복해서 수신, ReadMessage는 블로킹 함수
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// 클라이언트가 끊었을때 websocket: close 1000 (normal)
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
			// 타입에 따라 적절한 에러 처리 return도 가능.
			continue
		}
	}

	log.Println("Notificator service websocket close. ")
}
