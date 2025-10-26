package adapter

import "auth/internal/application/usecase/input"

func MakeUserAuthRegisterInput(id string, salt string, authHash string, userHash string) input.UserAuthRegisterInput {
	return input.UserAuthRegisterInput{
		Id:       id,
		Salt:     salt,
		UserHash: userHash,
		AuthHash: authHash,
	}
}
