package job

import (
	"log"
)

type ChatLineJob struct {
	// 이 작업이 사용할 Usecase 또는 필요한 데이터 필드
}

func NewChatLineJob() Job {

	return &ChatLineJob{}
}

// Job 인터페이스의 Execute() 메서드를 구현합니다.
func (j *ChatLineJob) Execute() error {
	// 실제 DB 저장 로직 (가장 무거운 작업)을 여기서 호출합니다.
	log.Println("채팅 라인 저장 작업 시작...")
	return nil
}
