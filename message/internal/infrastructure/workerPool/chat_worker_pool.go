package workerPool

import (
	"log"
	"message/internal/domain/chat/job"
	"message/internal/domain/chat/repository"
)

type chatWorkerPool struct {
	jobs       chan *job.ChatLineJob
	count      int
	repository repository.ChatRepository
}

type ChatWorkerPool interface {
	AddTask(j *job.ChatLineJob)
	Init()
	Stop() // 리소스 정리용 (선택적이지만 권장됨)
}

// 인터페이스 타입으로 반환
func NewChatWorkerPool(count int, repository repository.ChatRepository) ChatWorkerPool {

	// NewChatWorkerPool 함수가 구조체 포인터(*chatWorkerPool)를 생성한 후,
	// 이를 인터페이스 타입(ChatWorkerPool)으로 반환하는 방식은 Go에서 Factory 함수를 구현하는 표준적인 방법입니다.
	return &chatWorkerPool{
		jobs:       make(chan *job.ChatLineJob, count),
		count:      count,
		repository: repository,
	}
}

func (p *chatWorkerPool) AddTask(j *job.ChatLineJob) {
	// 실제 jobs 채널에 job을 넣는 로직
	p.jobs <- j
}

func (p *chatWorkerPool) Init() {
	// 워커 고루틴을 시작하는 로직
	for i := 0; i < p.count; i++ {
		go p.worker(i)
	}
}

// 워커 풀에서 Stop()의 주요 역할은 새로운 Job의 수신을 중단하고,
// 대기 중인 Job을 모두 처리한 후, 워커 고루틴들을 정상적으로 종료시키는 것입니다.
func (p *chatWorkerPool) Stop() {
	// 1. jobs 채널을 닫아 AddTask 호출을 막고,
	// 2. worker 고루틴들의 range 루프가 종료되도록 신호를 보냅니다.
	close(p.jobs) // for job := range p.jobs 루프를 종료
	log.Println("ChatWorkerPool Stopped. All jobs channel closed.")

	// NOTE: 실제 운영 환경에서는 모든 워커 고루틴이 종료될 때까지 기다리는 (예: sync.WaitGroup 사용) 로직이 추가되어야 합니다.
	// 여기서는 간단히 채널 close만 처리합니다.
}

// worker 메서드는 풀의 내부 구현 상세이며, 외부(Usecase)에서 호출할 필요가 없으므로 인터페이스에 포함되면 안 됩니다.
func (p *chatWorkerPool) worker(id int) {
	for job := range p.jobs {
		log.Println("data 수신 worker id ", id)
		job.Execute(p.repository) // DB 처리 로직 수행 호출
	}
}
