package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
)

func MakeUserAuthChallengeInput(id string, uuid string) input.UserAuthChallengeInput {
	return input.UserAuthChallengeInput{
		Id:   id,
		Uuid: uuid,
	}
}

func MakeUserAuthChallengeOutput(c string, s string) output.UserAuthChallengeOutput {
	return output.UserAuthChallengeOutput{
		Challenge: c,
		Salt:      s,
	}
}
