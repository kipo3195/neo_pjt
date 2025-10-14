package input

type GenerateAuthTokenInput struct {
	Id   string
	Uuid string
}

func MakeGenerateAuthTokenInput(id string, uuid string) GenerateAuthTokenInput {
	return GenerateAuthTokenInput{
		Id:   id,
		Uuid: uuid,
	}
}
