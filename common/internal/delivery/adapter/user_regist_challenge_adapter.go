package adapter

import "common/internal/application/usecase/input"

func MakeUserRegisterChallengeInput(id string, salt string) input.UserRegisterChallengeInput {

	return input.UserRegisterChallengeInput{
		Id:   id,
		Salt: salt,
	}
}
