package adapter

import "common/internal/application/usecase/input"

func MakeUserRegisterInput(id string, salt string, fv string) input.UserRegisterInput {

	return input.UserRegisterInput{
		Id:   id,
		Salt: salt,
		Fv:   fv,
	}
}
