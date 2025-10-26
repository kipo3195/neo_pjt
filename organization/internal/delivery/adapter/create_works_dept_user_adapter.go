package adapter

import "org/internal/application/usecase/input"

func MakeCreateWorksDeptUserInput(org string) input.CreateWorksDeptUserInput {
	return input.CreateWorksDeptUserInput{
		Org: org,
	}
}
