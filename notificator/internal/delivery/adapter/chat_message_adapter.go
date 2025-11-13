package adapter

import "notificator/internal/application/usecase/input"

func MakeChatMessageInput(t string, sendUserhash string, contents string, destUserHash []string) input.ChatMessageInput {
	return input.ChatMessageInput{
		Type:         t,
		SendUserHash: sendUserhash,
		Contents:     contents,
		DestUserHash: destUserHash,
	}
}
