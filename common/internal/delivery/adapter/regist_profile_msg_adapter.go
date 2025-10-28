package adapter

import "common/internal/application/usecase/input"

func MakeRegistProfileMsgInput(userId string, msg string) input.RegistProfileMsgInput {
	return input.RegistProfileMsgInput{
		UserId: userId,
		Msg:    msg,
	}
}
