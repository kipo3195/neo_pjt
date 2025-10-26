package adapter

import "auth/internal/application/usecase/input"

func MakeUserAuthInput(id string, fv string, uuid string) input.UserAuthInput {
	return input.UserAuthInput{
		Id:   id,
		Fv:   fv,
		Uuid: uuid,
	}
}
