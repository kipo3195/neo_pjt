package adapter

import "auth/internal/application/usecase/input"

func MakeGenerateAuthTokenInput(id string, uuid string, force bool) input.GenerateAuthTokenInput {
	return input.GenerateAuthTokenInput{
		Id:    id,
		Uuid:  uuid,
		Force: force,
	}
}
