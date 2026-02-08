package workerPool

import (
	"context"
	"hash/fnv"
	"log"
	"notificator/internal/consts"
	"notificator/internal/domain/chat/entity"
	"notificator/internal/domain/chat/job"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/external/ws/dto"
	response "notificator/pkg/dto"
	"sync"
	"time"
)

type chatCountWorkerPool struct {
	workerChans []chan *job.ChatCountJob
	numWorkers  int

	messageSender port.MessageSender
	wg            sync.WaitGroup // sync.WaitGroup은 Go 언어에서 여러 고루틴(Goroutine)이 모두 종료될 때까지 메인 흐름을 대기시키기 위해 사용하는 동기화 도구입니다.
	mu            sync.Mutex     // 상태 보호를 위한 뮤텍스 추가
	isClosed      bool
}

type ChatCountWorkerPool interface {
	AddTask(userHash string, en entity.ChatCountMessageEntity)
	Init()
	Stop() // 리소스 정리용 (선택적이지만 권장됨)
}

func NewChatCountWorkerPool(count int, messageSender port.MessageSender) *chatCountWorkerPool {

	// 각각의 워커가 자신만의 채널을 갖도록
	chans := make([]chan *job.ChatCountJob, count)
	for i := 0; i < count; i++ {
		chans[i] = make(chan *job.ChatCountJob, 100) // 적절한 버퍼 크기
	}

	return &chatCountWorkerPool{
		workerChans:   chans,
		numWorkers:    count,
		messageSender: messageSender,
	}
}

func (p *chatCountWorkerPool) Init() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		// 각 워커에게 고유 ID와 전용 채널을 할당하여 실행
		go p.worker(i, p.workerChans[i])
	}
	log.Printf("[chatCountWorkerPool] %d worker init.", p.numWorkers)
}

func (p *chatCountWorkerPool) worker(id int, ch chan *job.ChatCountJob) {
	defer p.wg.Done()

	// 워커마다 독립적인 맵을 생성합니다. (Lock-Free의 핵심)
	pendingMap := make(map[string]*entity.ChatCountJobEntity)

	for j := range ch {
		log.Printf("[chat count Worker %d] recv user : %s, isFlush : %v", id, j.UserHash, j.IsFlush)

		if j.IsFlush {
			// 타이머에 의한 전송 처리 로직
			p.handleFlush(j.UserHash, pendingMap, j.En)
		} else {
			// 신규 메시지 누적 및 디바운싱 로직
			p.handleMessage(ch, j, pendingMap)
		}
	}
}

func (p *chatCountWorkerPool) AddTask(userHash string, en entity.ChatCountMessageEntity) {

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
	case p.workerChans[workerIdx] <- &job.ChatCountJob{ // select 문 안의 case에 채널 전송 연산을 넣으면, **"해당 채널에 데이터를 즉시 보낼 수 있는 공간(Buffer)이 있는가?
		UserHash: userHash,
		En:       en,
		IsFlush:  false}:
	default:
		// 채널이 꽉 찼을 때의 전략 (에러 로그 또는 드롭)
		log.Printf("Chat Count Worker channel %d is full, dropping task for %s", workerIdx, userHash)
	}
}

func (p *chatCountWorkerPool) getWorkerIdx(userHash string) int {
	h := fnv.New32a()
	h.Write([]byte(userHash))
	// 32비트 해시 결과값을 워커 개수로 나눈 나머지(Modulo)를 구합니다.
	return int(h.Sum32()) % p.numWorkers
}

func (p *chatCountWorkerPool) handleMessage(ch chan *job.ChatCountJob, j *job.ChatCountJob, pendingMap map[string]*entity.ChatCountJobEntity) {

	item, exists := pendingMap[j.UserHash]

	// 읽음 처리면 바로 flush
	if j.En.EventType == consts.READ {

		log.Println("[handleMessage] chat read userHash :", j.UserHash)
		if exists && item.Timer != nil {
			item.Timer.Stop() // 대기 중인 unread 타이머 중단
		}

		delete(pendingMap, j.UserHash) // 버퍼 비우기

		// infra dto 형태로 파싱
		// 원본 엔티티의 델타를 버퍼링된 최종 합계로 교체
		chatCountData := dto.ChatCountDataDto{
			RoomKey:  j.En.ChatCountData.RoomKey,
			RoomType: j.En.ChatCountData.RoomType,
			Delta:    j.En.ChatCountData.Delta,
		}

		chatCountDataResponse := dto.ChatCountDataResponse{
			ChatCountData: chatCountData,
		}

		out := response.WsResponseDTO[dto.ChatCountDataResponse]{
			Type:      j.En.Type,
			EventType: j.En.EventType,
			Data:      chatCountDataResponse,
		}

		// 누적된 delta 무시하고 read 이벤트 즉시 전송 (보통 read는 delta가 0이거나 특정 고정값임)
		p.messageSender.SendToClient(j.UserHash, out)
		return
	}

	if !exists {

		item = &entity.ChatCountJobEntity{
			UserHash: j.UserHash,
			Delta:    j.En.ChatCountData.Delta,
			Count:    1,
		}
		pendingMap[j.UserHash] = item

		// 클로저 이슈 방지를 위한 값 복사
		targetUserHash := j.UserHash

		item.Timer = time.AfterFunc(500*time.Millisecond, func() {
			// 중요: 워커 자신의 채널로 Flush 작업을 다시 던집니다.
			ch <- &job.ChatCountJob{
				UserHash: targetUserHash,
				En:       j.En,
				IsFlush:  true,
			}
		})
	} else {
		// 기존 데이터에 Delta 누적
		item.Delta += j.En.ChatCountData.Delta
		item.Count++

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
	// 500ms 동안 쌓인 최종 Delta 값을 가져옵니다.
	finalDelta := item.Delta

	// 3. 버퍼에서 제거 (메모리 관리 및 다음 이벤트를 위한 초기화)
	delete(pendingMap, userHash)

	// 4. 전송용 DTO 구성
	// (여기서 r.En은 처음 타이머를 발동시킨 시점의 엔티티 정보를 활용하거나
	// Flush 전용 엔티티 구조를 설계하여 사용합니다.)
	chatCountData := dto.ChatCountDataDto{
		RoomKey:  en.ChatCountData.RoomKey, // item 엔티티에 Room 정보를 미리 담아두는 것이 좋습니다.
		RoomType: en.ChatCountData.RoomType,
		Delta:    finalDelta,
	}

	chatCountDataResponse := dto.ChatCountDataResponse{
		ChatCountData: chatCountData,
	}

	out := response.WsResponseDTO[dto.ChatCountDataResponse]{
		Type:      consts.CHATUNREAD,
		EventType: consts.UNREAD,
		Data:      chatCountDataResponse,
	}

	// 5. 클라이언트 전송 (Non-blocking 채널 전송이므로 워커가 멈추지 않음)
	p.messageSender.SendToClient(userHash, out)

	log.Printf("[Flush] User: %s, Final Delta: %d 전송 완료", userHash, finalDelta)
}

// 만들어 놓고 사용하지 않게 되있으므로 이후 호출 부분 추가하기
func (p *chatCountWorkerPool) Stop() {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return
	}
	p.isClosed = true // 1. 새로운 AddTask 호출 차단
	p.mu.Unlock()
	// close(ch)나 p.wg.Wait()을 Lock() 내부에 넣게 되면,
	// 워커들이 종료되면서 혹시라도 뮤텍스에 접근해야 하는 상황이 생길 때
	// 데드락이 발생할 위험이 있습니다

	log.Println("[chatWorkerPool] Stopping all workers...")

	// 2. 모든 워커 채널을 순회하며 닫음
	for _, ch := range p.workerChans {
		close(ch)
	}

	// 2. 타임아웃 컨텍스트 생성
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. WaitGroup 완료를 감지할 채널
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	// 4. 완료되었거나 시간이 다 되었거나
	select {
	case <-done:
		log.Println("[chatCountWorkerPool] All workers stopped gracefully.")
	case <-ctx.Done():
		log.Printf("[chatCountWorkerPool] Stop timed out after %v. Forcing shutdown.", 30)
		// 여기서 필요하다면 강제 종료를 위한 추가 로직 수행
	}

	log.Println("[chatWorkchatCountWorkerPoolerPool] All workers finished safely. Pool stopped.")
}
