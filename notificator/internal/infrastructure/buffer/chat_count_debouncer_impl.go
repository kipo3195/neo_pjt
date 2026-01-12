package buffer

import (
	"log"
	"notificator/internal/consts"
	"notificator/internal/domain/chat/entity"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/dto"
	"sync"
	"time"
)

type chatCountDebouncerImpl struct {
	sender  port.MessageSender
	mu      sync.Mutex
	buffers map[string]*pendingData
}

type pendingData struct {
	userHash string
	timer    *time.Timer
	count    int
	delta    int
}

func NewChatCountDebouncer(sender port.MessageSender) port.ChatCountDebouncer {
	return &chatCountDebouncerImpl{
		sender:  sender,
		buffers: make(map[string]*pendingData),
	}
}

func (r *chatCountDebouncerImpl) AddChatCount(userHash string, en *entity.ChatCountMessageEntity) error {

	r.mu.Lock()
	// 로직 종료 시 Lock 해제는 하단에서 명시적으로 하거나 상황에 따라 defer를 사용합니다.
	// 여기서는 가독성을 위해 defer를 유지하되, read 발생 시 로직을 분리합니다.

	item, exists := r.buffers[userHash]

	// [추가 정책] 만약 이번에 들어온 이벤트가 "read"라면?
	if en.EventType == consts.READ {
		log.Println("[AddChatCount] chat read userHash :", userHash)
		if exists && item.timer != nil {
			item.timer.Stop() // 대기 중인 unread 타이머 중단
		}
		delete(r.buffers, userHash) // 버퍼 비우기
		r.mu.Unlock()               // 즉시 전송을 위해 락 해제

		// infra dto 형태로 파싱
		// 원본 엔티티의 델타를 버퍼링된 최종 합계로 교체
		chatCountData := dto.ChatCountDataDto{
			RoomKey:  en.ChatCountData.RoomKey,
			RoomType: en.ChatCountData.RoomType,
			Delta:    en.ChatCountData.Delta,
		}

		chatCountMessage := dto.ChatCountMessageResponse{
			Type:          en.Type,
			EventType:     en.EventType,
			ChatCountData: chatCountData,
		}

		// 누적된 delta 무시하고 read 이벤트 즉시 전송 (보통 read는 delta가 0이거나 특정 고정값임)
		return r.sender.SendToClient(userHash, chatCountMessage)
	}

	// --- 기존 unread 디바운싱 로직 ---
	if !exists {
		// ... (1. 처음 들어온 unread 로직 - 기존과 동일) ...
		item = &pendingData{
			userHash: userHash,
			delta:    en.ChatCountData.Delta,
			count:    1,
		}
		r.buffers[userHash] = item
		item.timer = time.AfterFunc(500*time.Millisecond, func() {
			// 주의: Flush 내부에서 item.delta를 사용하므로 인자로 넘길 필요가 적어짐
			log.Println("[AddChatCount] chat unread 500ms Flush userHash :", userHash)
			r.Flush(userHash, en)
		})
		r.mu.Unlock()
	} else {
		// 2. 이미 unread 데이터가 쌓이고 있는 경우
		item.delta += en.ChatCountData.Delta
		item.count++

		if item.count >= 3 {
			if item.timer != nil {
				item.timer.Stop()
			}
			r.mu.Unlock() // Flush 전 락 해제
			log.Println("[AddChatCount] chat unread 3 count Flush userHash :", userHash)
			go r.Flush(userHash, en)
		} else {
			r.mu.Unlock()
		}
	}
	return nil
}

func (r *chatCountDebouncerImpl) Flush(userHash string, en *entity.ChatCountMessageEntity) {
	r.mu.Lock()
	item, exists := r.buffers[userHash]
	if !exists {
		r.mu.Unlock()
		return
	}

	// 전송할 최종 델타값 복사
	finalDelta := item.delta
	delete(r.buffers, userHash)
	r.mu.Unlock()

	// infra dto 형태로 파싱
	// 원본 엔티티의 델타를 버퍼링된 최종 합계로 교체
	chatCountData := dto.ChatCountDataDto{
		RoomKey:  en.ChatCountData.RoomKey,
		RoomType: en.ChatCountData.RoomType,
		Delta:    finalDelta,
	}

	chatCountMessage := dto.ChatCountMessageResponse{
		Type:          en.Type,
		EventType:     en.EventType,
		ChatCountData: chatCountData,
	}

	// 최종 전송
	r.sender.SendToClient(userHash, chatCountMessage)
}
