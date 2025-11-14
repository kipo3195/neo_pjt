package adapter

import (
	"user/internal/application/usecase/input"
)

func MakeRegistProfileMsgInput(userHash string, profileMsg string) input.RegistProfileMsgInput {
	return input.RegistProfileMsgInput{
		UserHash:   userHash,
		ProfileMsg: profileMsg,
	}
}
