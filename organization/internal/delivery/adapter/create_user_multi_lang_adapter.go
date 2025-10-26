package adapter

import "org/internal/application/usecase/input"

func MakeUserMultiLangInput(keyword string) input.CreateUserMultiLangInput {
	return input.CreateUserMultiLangInput{
		Keyword: keyword,
	}
}
