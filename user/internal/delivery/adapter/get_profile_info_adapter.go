package adapter

import (
	"user/internal/application/usecase/input"
	"user/internal/delivery/dto/userInfoService"
)

func MakeGetProfileInfoInput(users []userInfoService.UserInfoServiceDto) []input.GetUserDetailInfoInput {
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
