package adapter

import (
	"admin/internal/application/usecase/input"
	"admin/internal/application/usecase/output"
)

func MakeUserAuthInput(out output.RegistServiceUserOutput) input.UserAuthRegisterInput {

	in := make([]input.ServiceUserInput, 0)

	for _, su := range out.ServiceUser {

		temp := input.ServiceUserInput{
			UserHash: su.UserHash,
			UserId:   su.UserId,
			UserAuth: su.UserAuth,
			Salt:     su.Salt,
		}

		in = append(in, temp)

	}

	return input.UserAuthRegisterInput{
		ServiceUser: in,
	}
}
