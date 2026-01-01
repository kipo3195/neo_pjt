package workerPool

import (
	"log"
	"message/internal/domain/chat/job"
	"message/internal/domain/chat/repository"
	"sync"
)

type chatWorkerPool struct {
	jobs       chan *job.ChatLineJob
	count      int
	wg         sync.WaitGroup // sync.WaitGroup은 Go 언어에서 여러 고루틴(Goroutine)이 모두 종료될 때까지 메인 흐름을 대기시키기 위해 사용하는 동기화 도구입니다.
	mu         sync.Mutex     // 상태 보호를 위한 뮤텍스 추가
	repository repository.ChatRepository
	isClosed   bool
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

	p.mu.Lock()         // 잠금 시작
	defer p.mu.Unlock() // 함수가 끝날 때 잠금 해제

	if p.isClosed {
		log.Println("ChatWorkerPool Stopped.")
		return
	}

	// 실제 jobs 채널에 job을 넣는 로직
	p.jobs <- j
}

func (p *chatWorkerPool) Init() {
	// 워커 고루틴을 시작하는 로직
	// 워커 고루틴을 시작할때 wg.Add를 호출 ---- 1
	for i := 0; i < p.count; i++ {
		p.wg.Add(1) // 1을 더하기, Add는 반드시 고루틴 밖에서 호출하는 것이 안전합니다. (고루틴 안에서 호출하면 Wait이 먼저 실행될 위험이 있습니다.)
		go p.worker(i)
	}
}

// 워커 풀에서 Stop()의 주요 역할은 새로운 Job의 수신을 중단하고,
// 대기 중인 Job을 모두 처리한 후, 워커 고루틴들을 정상적으로 종료시키는 것입니다.
func (p *chatWorkerPool) Stop() {
	p.mu.Lock() // 잠금을 걸어서 AddTask가 동시에 실행되지 못하게 함
	if p.isClosed {
		// p.jobs가 이미 닫힌 상태면 panic이 발생 할 수 있으므로
		p.mu.Unlock()
		return
	}
	p.isClosed = true // 상태 변경
	// 1. jobs 채널을 닫아 AddTask 호출을 막고,
	// 2. worker 고루틴들의 range 루프가 종료되도록 신호를 보냅니다.
	close(p.jobs) // for job := range p.jobs 루프를 종료
	p.mu.Unlock() // 잠금 해제 (이제부터 AddTask는 위에서 차단됨)

	log.Println("ChatWorkerPool Stopped. All jobs channel closed.")

	// NOTE: 실제 운영 환경에서는 모든 워커 고루틴이 종료될 때까지 기다리는 (예: sync.WaitGroup 사용) 로직이 추가되어야 합니다.
	p.wg.Wait()
	log.Println("ChatWorkerPool Stopped. All workers finished safely.")
}

// worker 메서드는 풀의 내부 구현 상세이며, 외부(Usecase)에서 호출할 필요가 없으므로 인터페이스에 포함되면 안 됩니다.
func (p *chatWorkerPool) worker(id int) {
	// 종료될 때 wg.Done을 호출합니다. ---- 2
	defer p.wg.Done()
	for job := range p.jobs {
		log.Println("data 수신 worker id ", id)
		job.Execute(p.repository) // DB 처리 로직 수행 호출
	}
}
