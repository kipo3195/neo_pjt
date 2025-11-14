package adapter

import (
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
	"user/internal/domain/profile/entity"
)

func MakeGetProfileMsgInput(userHashs []string) input.GetProfileMsgInput {
	return input.GetProfileMsgInput{
		UserHashs: userHashs,
	}
}

func MakeGetProfileMsgOutput(result []entity.GetProfileMsgResultEntity) output.GetProfileMsgOutput {

	// nil slice → empty slice 로 보정
	if result == nil {
		result = []entity.GetProfileMsgResultEntity{}
	}

	return output.GetProfileMsgOutput{
		ProfileMsg: result,
	}
}
