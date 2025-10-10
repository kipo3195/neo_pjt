package input

type CreateUserMultiLangInput struct {
	Keyword string
}

func MakeUserMultiLangInput(keyword string) CreateUserMultiLangInput {
	return CreateUserMultiLangInput{
		Keyword: keyword,
	}
}
