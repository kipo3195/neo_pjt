package adapter

import "message/internal/application/usecase/input"

func MakeSendChatInput(lineKey string, contents string, destids []string) input.SendChatInput {
	return input.SendChatInput{
		LineKey:  lineKey,
		Contents: contents,
		DestIds:  destids,
	}
}
