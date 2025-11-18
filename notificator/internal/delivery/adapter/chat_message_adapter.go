package adapter

import "notificator/internal/application/usecase/input"

func MakeChatMessageInput(t string, sendUserhash string, contents string, lineKey string, recvUserHash []string) input.ChatMessageInput {
	return input.ChatMessageInput{
		Type:         t,
		SendUserHash: sendUserhash,
		Contents:     contents,
		LineKey:      lineKey,
		RecvUserHash: recvUserHash,
	}
}
