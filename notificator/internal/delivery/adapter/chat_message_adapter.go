package adapter

import "notificator/internal/application/usecase/input"

func MakeChatMessageInput(t string, sendUserhash string, contents string, lineKey string, destUserHash []string) input.ChatMessageInput {
	return input.ChatMessageInput{
		Type:         t,
		SendUserHash: sendUserhash,
		Contents:     contents,
		LineKey:      lineKey,
		DestUserHash: destUserHash,
	}
}
