package service

import (
	"context"
	"fmt"
	"log"
	"message/internal/application/usecase"
	"message/internal/application/usecase/input"
	"message/internal/consts"
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

	if in.TotalCount <= 0 || in.PerSecond <= 0 {
		return 0, consts.ErrBulkSendChatCountInvalid
	}

	targets, err := s.Chat.GetBulkSendChatTargets(ctx)
	if err != nil {
		return 0, err
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

	sentCount := 0
	for i := 0; i < in.TotalCount; i++ {
		select {
		case <-ctx.Done():
			return sentCount, ctx.Err()
		case <-ticker.C:
		}

		target := targets[i%len(targets)]
		if len(target.MemberHash) == 0 {
			continue
		}

		lineKey, sendDate, err := s.LineKey.GetLineKey(ctx)
		if err != nil {
			return sentCount, err
		}

		sendUserHash := target.MemberHash[i%len(target.MemberHash)]
		sendInput := input.SendChatInput{
			ChatRoom: input.ChatRoomInput{
				RoomKey:    target.RoomKey,
				RoomType:   target.RoomType,
				SecretFlag: target.SecretFlag,
			},
			ChatLine: input.ChatLineInput{
				Cmd:          cmd,
				Contents:     fmt.Sprintf("%s %06d", contents, i+1),
				LineKey:      lineKey,
				SendUserHash: sendUserHash,
				SendDate:     sendDate,
			},
			EventType:   eventType,
			ChatSession: chatSession,
		}

		if _, err := s.Chat.SendChat(ctx, sendInput); err != nil {
			log.Println("[SendBulkChat] send failed. index :", i+1, "roomKey :", target.RoomKey, "sendUserHash :", sendUserHash, "err :", err)
			return sentCount, err
		}
		sentCount++
	}

	return sentCount, nil
}
