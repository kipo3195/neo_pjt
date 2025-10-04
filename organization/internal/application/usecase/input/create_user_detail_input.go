package input

type CreateUserDetailInput struct {
	Keyword string
	Type    string
}

func MakeCreateUserDetailInput(k string, t string) CreateUserDetailInput {
	return CreateUserDetailInput{
		Keyword: k,
		Type:    t,
	}
}
