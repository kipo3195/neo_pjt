package input

type GenerateAuthTokenInput struct {
	Id    string
	Uuid  string
	Force bool
}

func MakeGenerateAuthTokenInput(id string, uuid string, force bool) GenerateAuthTokenInput {
	return GenerateAuthTokenInput{
		Id:    id,
		Uuid:  uuid,
		Force: force,
	}
}
