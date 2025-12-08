package workerPool

import (
	"message/internal/domain/chat/job"
	"message/internal/domain/chat/pool"
)

type chatWorkerPool struct {
	jobs  chan job.ChatLineJob // Job 인터페이스 타입의 채널
	count int
}

// Pool 구조체 인스턴스를 직접 사용하기 위한 처리
func NewChatWorkerPool(count int) pool.ChatPool {
	// Go 서비스 컴포넌트(특히 상태를 가지고 관리하는 풀이나 서비스)는 포인터를 반환하여 일관성을 유지하고 불필요한 값 복사를 막습니다.
	return &chatWorkerPool{
		count: count,
	}
}

func (p *chatWorkerPool) AddTask(j *job.ChatLineJob) {
	// (참고: 우선순위 처리는 별도의 복잡한 로직이 필요하며, 여기서는 단순화하여 처리합니다.)

	// 채널에 Job 인터페이스 객체를 넣습니다.
	p.jobs <- j
}

// // 워커 고루틴의 핵심 로직
// func worker(id int, jobs <-chan job.Job) {
// 	for job := range jobs {
// 		// 동적으로 job의 Execute() 메서드를 호출합니다.
// 		err := job.Execute()
// 		if err != nil {
// 			// 에러 처리 로직
// 		}
// 	}
// }
