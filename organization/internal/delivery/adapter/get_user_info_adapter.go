package adapter

import (
	"org/internal/application/usecase/input"
)

func MakeGetUserInfoInput(userIds []string) input.GetUserInfoInput {

	return input.GetUserInfoInput{
		UserIds: userIds,
	}
}
