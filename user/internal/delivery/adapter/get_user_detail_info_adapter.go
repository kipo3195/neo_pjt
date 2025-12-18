package adapter

import (
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
	"user/internal/delivery/dto/userInfoService"
	"user/internal/domain/userDetail/entity"
)

func MakeGetUserDetailInfoInput(users []userInfoService.UserInfoServiceDto) []input.GetUserDetailInfoInput {

	reqUsers := make([]input.GetUserDetailInfoInput, 0)

	for _, u := range users {

		temp := input.GetUserDetailInfoInput{
			UserHash:   u.UserHash,
			UpdateHash: u.UpdateHash,
		}

		reqUsers = append(reqUsers, temp)
	}

	return reqUsers
}

func MakeGetUserDetailInfoOutput(userInfos []entity.UserDetailInfoEntity) output.GetUserDetailInfoOutput {
	return output.GetUserDetailInfoOutput{
		UserInfos: userInfos,
	}
}
