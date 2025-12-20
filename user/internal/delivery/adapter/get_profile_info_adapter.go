package adapter

import (
	"user/internal/application/usecase/input"
	"user/internal/delivery/dto/userInfoService"
)

func MakeGetProfileInfoInput(users []userInfoService.UserInfoServiceDto) []input.GetUserProfileInfoInput {
	reqUsers := make([]input.GetUserProfileInfoInput, 0)

	for _, u := range users {

		temp := input.GetUserProfileInfoInput{
			UserHash:       u.UserHash,
			ProfileVersion: u.ProfileVersion,
		}

		reqUsers = append(reqUsers, temp)
	}

	return reqUsers

}
