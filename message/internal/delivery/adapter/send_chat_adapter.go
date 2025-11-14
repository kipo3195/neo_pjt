package adapter

import "message/internal/application/usecase/input"

func MakeSendChatInput(sendUserHash string, lineKey string, contents string, destUsers []string) input.SendChatInput {
	return input.SendChatInput{
		SendUserHash: sendUserHash,
		LineKey:      lineKey,
		Contents:     contents,
		DestUsers:    destUsers,
	}
}
