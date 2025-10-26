package adapter

import "common/internal/application/usecase/input"

func MakeGetProfileImgInput(userId string) input.GetProfileImgInput {
	return input.GetProfileImgInput{
		UserId: userId,
	}
}
