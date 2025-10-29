package adapter

import "user/internal/application/usecase/input"

func MakeGetProfileImgInput(userId string) input.GetProfileImgInput {
	return input.GetProfileImgInput{
		UserId: userId,
	}
}
