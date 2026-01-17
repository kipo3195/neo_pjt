package adapter

import (
	"admin/internal/application/usecase/input"
	"admin/internal/application/usecase/output"
)

func MakePublishServiceUserInput(org string, out output.RegistServiceUserOutput) []input.PublishServiceUserInput {

	in := make([]input.PublishServiceUserInput, 0)

	for _, su := range out.ServiceUser {

		temp := input.PublishServiceUserInput{
			Org:      org,
			UserHash: su.UserHash,
			UserId:   su.UserId,
		}

		in = append(in, temp)

	}
	return in

}
