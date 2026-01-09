package job

import (
	"context"
	"log"
	"message/internal/consts"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/repository"
	"message/internal/util"

	"github.com/nats-io/nats.go"
)

// Job의 역할: 도메인 로직 수행자
// ChatLineJob의 Execute() 메서드는 최종적으로 DB 저장이라는 무거운 I/O 작업을 수행하지만,
// 이 작업의 '무엇을' 저장하고 '어떻게' 저장할지 결정하는 것은 비즈니스 로직의 영역입니다.
// ChatLineJob은 LineKey 같은 도메인 데이터를 캡슐화하고 있으며, 이를 처리하는 **행위(Execute)**를 정의합니다.
type ChatLineJob struct {
	SendChatEntity   entity.SendChatEntity
	ChatUnreadEntity entity.ChatUnreadEntity

	// 🎯 Context 필드 추가
	Ctx       context.Context
	Cancel    context.CancelFunc // 👈 Cancel 함수를 Job에 포함
	Connector *nats.Conn
}

// Execute 메서드는 워커 풀이 주입해준 Repository를 사용하여 작업을 수행합니다.
// LineKey는 Job 자체에 이미 포함되어 있으므로 인자로 받지 않습니다.
func (j *ChatLineJob) Execute(repository repository.ChatRepository) error {
	defer j.Cancel()

	// 실제 DB 저장 로직 (가장 무거운 작업)을 여기서 호출합니다.
	// 예시: LineKey를 사용하여 데이터를 찾거나 저장합니다.
	err := repository.SaveChatLine(j.Ctx, j.SendChatEntity)
	if err != nil {
		return err
	}

	log.Println("[SendChatUnread] send entity : ", j.ChatUnreadEntity)

	data, err := util.EntityMarshal(j.ChatUnreadEntity)
	if err != nil {
		log.Fatal(err)
		return err
	}

	/* 채팅 발송 Message Broker */
	err = j.Connector.Publish("chat.unread.broadcast", data)
	if err != nil {
		log.Fatal("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError
		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}

	return nil
}
