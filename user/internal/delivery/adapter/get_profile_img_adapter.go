package adapter

import "user/internal/application/usecase/input"

func MakeGetProfileImgInput(userHash string) input.GetProfileImgInput {
	return input.GetProfileImgInput{
		UserHash: userHash,
	}
}
