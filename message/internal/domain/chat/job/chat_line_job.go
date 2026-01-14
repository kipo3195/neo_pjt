package job

import (
	"context"
	"message/internal/domain/chat/entity"

	"github.com/nats-io/nats.go"
)

// Job의 역할: 도메인 로직 수행자
// ChatLineJob의 Execute() 메서드는 최종적으로 DB 저장이라는 무거운 I/O 작업을 수행하지만,
// 이 작업의 '무엇을' 저장하고 '어떻게' 저장할지 결정하는 것은 비즈니스 로직의 영역입니다.
// ChatLineJob은 LineKey 같은 도메인 데이터를 캡슐화하고 있으며, 이를 처리하는 **행위(Execute)**를 정의합니다.
type ChatLineJob struct {
	SendChatEntity       entity.SendChatEntity
	ChatCountEventEntity entity.ChatCountEventEntity

	// 🎯 Context 필드 추가
	Ctx       context.Context
	Cancel    context.CancelFunc // 👈 Cancel 함수를 Job에 포함
	Connector *nats.Conn
}
