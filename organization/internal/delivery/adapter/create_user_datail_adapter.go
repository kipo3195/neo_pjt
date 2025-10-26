package adapter

import "org/internal/application/usecase/input"

func MakeCreateUserDetailInput(k string, t string) input.CreateUserDetailInput {
	return input.CreateUserDetailInput{
		Keyword: k,
		Type:    t,
	}
}
