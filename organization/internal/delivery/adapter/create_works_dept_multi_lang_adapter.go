package adapter

import "org/internal/application/usecase/input"

func MakeCreateWorksDeptMultiLangInput(org string) input.CreateWorksDeptMultiLangInput {
	return input.CreateWorksDeptMultiLangInput{
		Org: org,
	}
}
