package input

type CreateWorksDeptMultiLangInput struct {
	Org string
}

func MakeCreateWorksDeptMultiLangInput(org string) CreateWorksDeptMultiLangInput {
	return CreateWorksDeptMultiLangInput{
		Org: org,
	}
}
