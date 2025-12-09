package job

import (
	"context"
	"log"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/repository"
)

// Job의 역할: 도메인 로직 수행자
// ChatLineJob의 Execute() 메서드는 최종적으로 DB 저장이라는 무거운 I/O 작업을 수행하지만,
// 이 작업의 '무엇을' 저장하고 '어떻게' 저장할지 결정하는 것은 비즈니스 로직의 영역입니다.
// ChatLineJob은 LineKey 같은 도메인 데이터를 캡슐화하고 있으며, 이를 처리하는 **행위(Execute)**를 정의합니다.
type ChatLineJob struct {
	SendChatEntity entity.SendChatEntity
	// 🎯 Context 필드 추가
	Ctx context.Context
}

// Execute 메서드는 워커 풀이 주입해준 Repository를 사용하여 작업을 수행합니다.
// LineKey는 Job 자체에 이미 포함되어 있으므로 인자로 받지 않습니다.
func (j *ChatLineJob) Execute(repository repository.ChatRepository) error {
	log.Println("채팅 라인 저장 작업 시작. LineKey:", j.SendChatEntity.ChatLineEntity.LineKey)

	// 실제 DB 저장 로직 (가장 무거운 작업)을 여기서 호출합니다.
	// 예시: LineKey를 사용하여 데이터를 찾거나 저장합니다.
	repository.SaveChatLine(j.Ctx, j.SendChatEntity)

	return nil
}
