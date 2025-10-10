package input

type CreateWorksDeptUserInput struct {
	Org       string
	UserHashs []string
}

func MakeCreateWorksDeptUserInput(org string) CreateWorksDeptUserInput {
	return CreateWorksDeptUserInput{
		Org: org,
	}
}
