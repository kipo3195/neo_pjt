package workerPool

import (
	"hash/fnv"
	"log"
	"notificator/internal/domain/chat/entity"
	"notificator/internal/domain/chat/job"
	"notificator/internal/domain/port"
	"sync"
	"time"
)

type chatReadDateWorkerPool struct {
	workerChans []chan *job.ChatReadDateJob
	numWorkers  int

	messageSender port.MessageSender
	wg            sync.WaitGroup // sync.WaitGroup은 Go 언어에서 여러 고루틴(Goroutine)이 모두 종료될 때까지 메인 흐름을 대기시키기 위해 사용하는 동기화 도구입니다.
	mu            sync.Mutex     // 상태 보호를 위한 뮤텍스 추가
	isClosed      bool
}

type ChatReadDateWorkerPool interface {
	AddTask(userHash string, en entity.ChatReadMessageEntity)
	Init()
	Stop() // 리소스 정리용 (선택적이지만 권장됨)
}

func NewChatReadDateWorkerPool(count int, messageSender port.MessageSender) *chatReadDateWorkerPool {
	// 각각의 워커가 자신만의 채널을 갖도록
	chans := make([]chan *job.ChatReadDateJob, count)
	for i := 0; i < count; i++ {
		chans[i] = make(chan *job.ChatReadDateJob, 100) // 적절한 버퍼 크기
	}

	return &chatReadDateWorkerPool{
		workerChans:   chans,
		numWorkers:    count,
		messageSender: messageSender,
	}
}

func (p *chatReadDateWorkerPool) Init() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		// 각 워커에게 고유 ID와 전용 채널을 할당하여 실행
		go p.worker(i, p.workerChans[i])
	}
	log.Printf("[chatReadDateWorkerPool] %d worker init.", p.numWorkers)
}

func (p *chatReadDateWorkerPool) worker(id int, ch chan *job.ChatReadDateJob) {
	defer p.wg.Done()

	// 워커마다 독립적인 맵을 생성합니다. (Lock-Free의 핵심)
	pendingMap := make(map[string]entity.ChatReadDateJobEntity)
	//id : roomKey : [{member, readDate}]

	for j := range ch {
		log.Printf("[chat readDate Worker %d] recv user : %s, isFlush : %v", id, j.UserHash, j.IsFlush)

		if j.IsFlush {
			// 타이머에 의한 전송 처리 로직
			p.handleFlush(j.UserHash, pendingMap, j.En)
		} else {
			// 신규 메시지 누적 및 디바운싱 로직
			p.handleMessage(ch, j, pendingMap)
		}
	}
}

func (p *chatReadDateWorkerPool) AddTask(userHash string, en entity.ChatReadMessageEntity) {

	// 빠르게 던지고 빠짐
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	// 1. 유저 해시 기반 인덱스 추출
	workerIdx := p.getWorkerIdx(userHash)

	// 2. 해당 워커 전용 채널에 투척
	select {
	case p.workerChans[workerIdx] <- &job.ChatReadDateJob{ // select 문 안의 case에 채널 전송 연산을 넣으면, **"해당 채널에 데이터를 즉시 보낼 수 있는 공간(Buffer)이 있는가?
		UserHash: userHash,
		En:       en,
		IsFlush:  false}:
	default:
		// 채널이 꽉 찼을 때의 전략 (에러 로그 또는 드롭)
		log.Printf("Chat Count Worker channel %d is full, dropping task for %s", workerIdx, userHash)
	}
}

func (p *chatReadDateWorkerPool) getWorkerIdx(userHash string) int {
	h := fnv.New32a()
	h.Write([]byte(userHash))
	// 32비트 해시 결과값을 워커 개수로 나눈 나머지(Modulo)를 구합니다.
	return int(h.Sum32()) % p.numWorkers
}

func (p *chatReadDateWorkerPool) handleMessage(ch chan *job.ChatReadDateJob, j *job.ChatReadDateJob, pendingMap map[string]entity.ChatReadDateJobEntity) {

	item, exists := pendingMap[j.UserHash]

	if !exists {

		// [{member, readDate}, {member, readDate}] 인데.. 결국 set의 형태가 되어야 하는듯 TODO
		readMember := make([]entity.ChatReadDateEntity, 0)
		readMember = append(readMember, entity.ChatReadDateEntity{
			MemberHash: j.En.ChatReadData.MemberHash,
			ReadDate:   j.En.ChatReadData.ReadDate,
		})

		// room : [{member, readDate}, {member, readDate}]
		readDateMap := make(map[string][]entity.ChatReadDateEntity)

		readDateMap[j.En.ChatReadData.RoomKey] = readMember

		item.RoomReadMap = readDateMap

		item.Timer = time.AfterFunc(500*time.Millisecond, func() {
			// 중요: 워커 자신의 채널로 Flush 작업을 다시 던집니다.
			ch <- &job.ChatReadDateJob{
				UserHash: j.UserHash,
				En:       j.En,
				IsFlush:  true,
			}
		})
	} else {

		// [선택 사항] 누적 횟수가 특정 횟수를 넘어서면 타이머 무시하고 즉시 Flush
		if item.Count >= 5 {
			if item.Timer != nil {
				item.Timer.Stop()
			}
			p.handleFlush(j.UserHash, pendingMap, j.En)
		}
	}
}

func (p *chatCountWorkerPool) handleFlush(userHash string, pendingMap map[string]*entity.ChatCountJobEntity, en entity.ChatCountMessageEntity) {
	// 1. 해당 유저의 펜딩 데이터가 있는지 확인
	item, exists := pendingMap[userHash]
	if !exists {
		return
	}

	// 2. 전송할 데이터 스냅샷 생성

	// 3. 버퍼에서 제거 (메모리 관리 및 다음 이벤트를 위한 초기화)
	delete(pendingMap, userHash)

	// 4. 전송용 DTO 구성
	// (여기서 r.En은 처음 타이머를 발동시킨 시점의 엔티티 정보를 활용하거나
	// Flush 전용 엔티티 구조를 설계하여 사용합니다.)

	//dto

	// res

	// out

	// 5. 클라이언트 전송 (Non-blocking 채널 전송이므로 워커가 멈추지 않음)
	p.messageSender.SendToClient(userHash, out)

	log.Printf("[Flush] User: %s, Final Delta: %d 전송 완료", userHash, finalDelta)
}
