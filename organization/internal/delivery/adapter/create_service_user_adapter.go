package adapter

import "org/internal/application/usecase/input"

func MakeCreateServiceUserInput(userCount int, keyword string) input.CreateServiceUserInput {
	return input.CreateServiceUserInput{
		UserCount: userCount,
		Keyword:   keyword,
	}
}
