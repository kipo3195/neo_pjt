package adapter

import "common/internal/application/usecase/input"

func MakeProfileUploadInput(profileImg *[]byte, profileImgSize int64, profileImgName string, userId string) input.ProfileImgInput {
	return input.ProfileImgInput{
		ProfileImg:     profileImg,
		ProfileImgSize: profileImgSize,
		ProfileImgName: profileImgName,
		UserId:         userId,
	}
}
