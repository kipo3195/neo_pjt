package adapter

import (
	"org/internal/application/usecase/input"
)

func MakeGetUserInfoInput(userHashs []string) input.GetUserInfoInput {

	return input.GetUserInfoInput{
		UserHashs: userHashs,
	}
}
