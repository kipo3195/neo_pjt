package adapter

import "org/internal/application/usecase/input"

func MakeMyInfoInput(myHash string) input.MyInfoInput {
	return input.MyInfoInput{
		MyHash: myHash,
	}
}
