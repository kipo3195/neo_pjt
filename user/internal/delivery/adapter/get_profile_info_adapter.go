package adapter

import "user/internal/application/usecase/input"

func MakeGetProfileInfoInput(userhashs []string) input.GetProfileInfoInput {
	return input.GetProfileInfoInput{
		UserHashs: userhashs,
	}

}
