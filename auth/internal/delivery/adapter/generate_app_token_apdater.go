package adapter

import "auth/internal/application/usecase/input"

func MakeGenerateAppTokenInput(uuid string) input.GenerateAppTokenInput {

	return input.GenerateAppTokenInput{
		Uuid: uuid,
	}
}
