package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/delivery/dto/userAuth"
)

func MakeUserAuthRegisterInput(dto []userAuth.UserAuthRegisterDto) []input.UserAuthRegisterInput {

	in := make([]input.UserAuthRegisterInput, 0)

	for _, d := range dto {

		temp := input.UserAuthRegisterInput{
			Id:       d.Id,
			Salt:     d.Salt,
			UserHash: d.UserHash,
			AuthHash: d.AuthHash,
		}
		in = append(in, temp)

	}
	return in
}
