package adapter

import "user/internal/application/usecase/input"

func MakeDeleteProfileImgInput(userHash string) input.DeleteProfileImgInput {
	return input.DeleteProfileImgInput{
		UserHash: userHash,
	}
}
