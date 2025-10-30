package adapter

import (
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
	"user/internal/domain/userDetail/entity"
)

func MakeGetUserDetailInfoInput(userhashs []string) input.GetUserDetailInfoInput {

	return input.GetUserDetailInfoInput{
		UserHashs: userhashs,
	}
}

func MakeGetUserDetailInfoOutput(userInfos []entity.UserDetailInfoEntity) output.GetUserDetailInfoOutput {
	return output.GetUserDetailInfoOutput{
		UserInfos: userInfos,
	}
}
