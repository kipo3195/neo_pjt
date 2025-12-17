package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/delivery/dto/userAuth"
)

func MakeUserAuthRegisterInput(dto []userAuth.UserAuthRegisterDto) []input.UserAuthRegisterInput {

	in := make([]input.UserAuthRegisterInput, 0)

	for _, d := range dto {

		temp := input.UserAuthRegisterInput{
			UserId:   d.UserId,
			Salt:     d.Salt,
			UserHash: d.UserHash,
			UserAuth: d.UserAuth,
		}
		in = append(in, temp)

	}
	return in
}
