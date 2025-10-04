package input

type CreateServiceUserInput struct {
	UserCount int
	Keyword   string
}

func MakeCreateServiceUserInput(userCount int, keyword string) CreateServiceUserInput {
	return CreateServiceUserInput{
		UserCount: userCount,
		Keyword:   keyword,
	}
}
