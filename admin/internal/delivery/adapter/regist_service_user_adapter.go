package adapter

import (
	"admin/internal/application/usecase/input"
	"admin/internal/application/usecase/output"
	"admin/internal/domain/serviceUser/entity"
)

func MakeRegistServiceUserInput(org string, userId []string, userAuth string) input.RegistServiceUserInput {
	return input.RegistServiceUserInput{
		Org:      org,
		UserId:   userId,
		UserAuth: userAuth,
	}
}

func MakeRegistServiceUserOutput(entity []entity.ServiceUserEntity) output.RegistServiceUserOutput {

	serviceUser := make([]output.ServiceUserOutput, 0)

	for _, e := range entity {

		value := output.ServiceUserOutput{
			UserId:   e.UserId,
			UserHash: e.UserHash,
		}

		serviceUser = append(serviceUser, value)
	}

	return output.RegistServiceUserOutput{
		ServiceUser: serviceUser,
	}
}
