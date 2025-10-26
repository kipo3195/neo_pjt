package adapter

import "common/internal/application/usecase/input"

func MakeDeleteProfileImgInput(userId string) input.DeleteProfileImgInput {
	return input.DeleteProfileImgInput{
		UserId: userId,
	}
}
