package adapter

import (
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
	"user/internal/delivery/dto/userInfoService"
	"user/internal/domain/userDetail/entity"
)

func MakeGetUserInfoInput(users []userInfoService.UserInfoServiceDto) ([]input.GetUserDetailInfoInput, []input.GetUserProfileInfoInput) {

	detail := make([]input.GetUserDetailInfoInput, 0)
	profile := make([]input.GetUserProfileInfoInput, 0)

	for _, u := range users {

		d := input.GetUserDetailInfoInput{
			UserHash:   u.UserHash,
			DetailHash: u.DetailHash,
		}

		p := input.GetUserProfileInfoInput{
			UserHash:    u.UserHash,
			ProfileHash: u.ProfileHash,
		}

		detail = append(detail, d)
		profile = append(profile, p)
	}

	return detail, profile
}

func MakeGetUserDetailInfoOutput(userInfos []entity.UserDetailInfoEntity) []output.UserDetailOutput {

	result := make([]output.UserDetailOutput, 0)

	for _, u := range userInfos {
		temp := output.UserDetailOutput{

			UserHash:     u.UserHash,
			UserEmail:    u.UserEmail,
			UserPhoneNum: u.UserPhoneNum,
		}
		result = append(result, temp)
	}

	return result
}
