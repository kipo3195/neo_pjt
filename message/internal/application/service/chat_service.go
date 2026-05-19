package service

import (
	"context"
	"fmt"
	"log"
	"message/internal/application/usecase"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"sync"
	"time"
)

type ChatService struct {
	Chat     usecase.ChatUsecase
	LineKey  usecase.LineKeyUsecase
	ChatRoom usecase.ChatRoomUsecase
}

func NewChatService(chat usecase.ChatUsecase, lineKey usecase.LineKeyUsecase, chatRoom usecase.ChatRoomUsecase) *ChatService {
	return &ChatService{
		Chat:     chat,
		LineKey:  lineKey,
		ChatRoom: chatRoom,
	}
}

func (s *ChatService) SendBulkChat(ctx context.Context, in input.BulkSendChatInput) (int, error) {

	if in.TotalCount <= 0 || in.PerSecond <= 0 || in.WorkerCount < 0 {
		return 0, consts.ErrBulkSendChatCountInvalid
	}

	// roomKey 조회
	targets, err := s.Chat.GetBulkSendChatTargets(ctx)
	if err != nil {
		return 0, err
	}

	workerCount := in.WorkerCount
	if workerCount == 0 {
		workerCount = 1
	}

	contents := in.Contents
	if contents == "" {
		contents = "bulk test chat"
	}

	cmd := in.Cmd
	if cmd == 0 {
		cmd = 1
	}

	eventType := in.EventType
	if eventType == "" {
		eventType = "C"
	}

	chatSession := in.ChatSession
	if chatSession == "" {
		chatSession = "bulk"
	}

	interval := time.Second / time.Duration(in.PerSecond)
	if interval <= 0 {
		interval = time.Nanosecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	type bulkChatJob struct {
		Index   int
		RoomKey string
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobCh := make(chan bulkChatJob)
	resultCh := make(chan error, workerCount)

	// 요청받은 workerCount 만큼 고루틴 생성, 채널의 데이터를 처리
	var wg sync.WaitGroup
	for workerID := 0; workerID < workerCount; workerID++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for job := range jobCh {
				select {
				case <-runCtx.Done():
					return
				case <-ticker.C:
				}

				lineKey, sendDate, err := s.LineKey.GetLineKey(runCtx)
				if err != nil {
					resultCh <- err
					cancel()
					return
				}

				sendInput := input.SendChatInput{
					ChatRoom: input.ChatRoomInput{
						RoomKey:    job.RoomKey,
						RoomType:   "N",
						SecretFlag: "N",
					},
					ChatLine: input.ChatLineInput{
						Cmd:          cmd,
						Contents:     fmt.Sprintf("%s %06d", contents, job.Index+1),
						LineKey:      lineKey,
						SendUserHash: "bulk",
						SendDate:     sendDate,
					},
					EventType:   eventType,
					ChatSession: chatSession,
				}

				if _, err := s.Chat.SendChat(runCtx, sendInput); err != nil {
					log.Println("[SendBulkChat] send failed. workerID :", workerID, "index :", job.Index+1, "roomKey :", job.RoomKey, "sendUserHash : bulk", "err :", err)
					resultCh <- err
					cancel()
					return
				}

				resultCh <- nil
			}
		}(workerID)
	}

	// 요청받은 채팅 데이터의 건수를 정확하게 하기 위해서 TotalCount로 순회
	// 하지만 방은 순서대로 하되, TotalCount까지.
	go func() {
		defer close(jobCh)

		for i := 0; i < in.TotalCount; i++ {
			target := targets[i%len(targets)]
			select {
			case <-runCtx.Done():
				return
			case jobCh <- bulkChatJob{
				Index:   i,
				RoomKey: target.RoomKey,
			}:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	sentCount := 0
	var firstErr error
	for err := range resultCh {
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		sentCount++
	}

	if firstErr != nil {
		return sentCount, firstErr
	}
	if err := ctx.Err(); err != nil && sentCount < in.TotalCount {
		return sentCount, err
	}

	return sentCount, nil
}
