package usecase

import (
	"log"
	"notificator/internal/domain/socketSender/entity"
	"notificator/internal/infrastructure/config"
	"notificator/internal/infrastructure/storage"
	"time"

	"github.com/gorilla/websocket"
)

type socketSenderUsecase struct {
	sendConnectionStorage storage.SendConnectionStorage
}

type SocketSenderUsecase interface {
	SaveConnection(conn *websocket.Conn, userHash string, websocketConfig config.WebsocketConnectionConfig)
	GetConnection(userHash string) *entity.SendConnectionEntity
}

func NewSocketSenderUsecase(sendConnectionStorage storage.SendConnectionStorage) SocketSenderUsecase {
	return &socketSenderUsecase{
		sendConnectionStorage: sendConnectionStorage,
	}
}

func (r *socketSenderUsecase) GetConnection(userHash string) *entity.SendConnectionEntity {
	return r.sendConnectionStorage.GetConnection(userHash)
}

func (r *socketSenderUsecase) SaveConnection(conn *websocket.Conn, userHash string, websocketConfig config.WebsocketConnectionConfig) {

	// 1. 설정값 정의
	var pingPeriod = time.Duration(websocketConfig.PingPeriod) * time.Second // Ping 발송 주기
	var writeWait = time.Duration(websocketConfig.WriteWait) * time.Second   // 쓰기 타임아웃

	// 쓰기용 채널 생성 (Write Channel)
	Ch := make(chan interface{}, 100) // 초당 3~4건의 데이터를 받는다고 가정했을때 25~30초 분량의 트래픽을 받아 낼 수 있음.
	entity := entity.MakeSendConnectionEntity(userHash, conn, Ch)
	r.sendConnectionStorage.PutConnection(userHash, entity)

	// Ping을 위한 티커 생성
	ticker := time.NewTicker(pingPeriod)

	// 함수 종료시 티커 중지 + 세션 삭제 + 소켓 종료
	defer func() {
		ticker.Stop()                                      // 티커 중지
		r.sendConnectionStorage.RemoveConnection(userHash) // 메모리 누수 및 블로킹 방지를 위한 로직 -> 이 userHash로 메시지를 보내려 해도 저장소에 없으므로...
		conn.Close()                                       // 소켓 연결 종료
	}()

	log.Printf("[SaveConnection] Start for user: %s", userHash)

	for {
		select {

		case message, ok := <-entity.Chan: // 일반 메시지 수신
			// 쓰기 타임아웃 설정
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// 외부에서 채널을 닫았을 경우 (정상 종료 시나리오)
				log.Printf("[SaveConnection] Channel closed for user: %s", userHash)
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteJSON(message); err != nil {
				log.Printf("[SaveConnection] WriteJSON error: %v", err)
				return // 에러 발생 시 고루틴 종료 (defer 실행)
			}

		case <-ticker.C: // 주기적인 Ping 발송 시간
			// 쓰기 타임아웃 설정
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			log.Println("[SaveConnection] Sending Ping to user:", userHash)
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[SaveConnection] Ping error for user %s: %v", userHash, err)
				return // 에러 발생 시 고루틴 종료 (defer 실행)
			}
		}
	}

	// 하위 로직은 for-select 문 안으로 이동하여 관리됩니다.
	/* 실제로 클라이언트 소켓에 write하는 로직 */
	// for message := range entity.Chan {
	// 	// 'ok' 검사가 필요 없습니다. 채널이 닫히면 for 루프는 자동으로 종료됩니다.

	// 	// 이 위치에서 conn.WriteJSON 호출 (단일 고루틴에서만 접근하므로 안전)
	// 	if err := entity.Conn.WriteJSON(message); err != nil {
	// 		// 소켓이 닫혔을때 (클라이언트에서 닫으면, err -> websocket: close sent) 채널을 통해 데이터를 수신하지만, write할때 close send가 발생하므로
	// 		//  고루틴 종료를 위해 break를 통해 종료 하여 루프를 빠져나가게 한다.
	// 		log.Printf("[SaveConnection] userHash :%s, error : %s \n", userHash, err)
	// 		break
	// 	}
	// }

	// 만약 쓰기 고루틴에서 에러가 발생하여 소켓을 닫을 일이 있다면 여기서 conn을 닫는것도 방법임.
}
