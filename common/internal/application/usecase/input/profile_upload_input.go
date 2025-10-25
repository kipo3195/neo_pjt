package input

type ProfileImgInput struct {
	ProfileImg     *[]byte
	ProfileImgSize int64
	ProfileImgName string
	UserId         string
}

func MakeProfileUploadInput(profileImg *[]byte, profileImgSize int64, profileImgName string, userId string) ProfileImgInput {
	return ProfileImgInput{
		ProfileImg:     profileImg,
		ProfileImgSize: profileImgSize,
		ProfileImgName: profileImgName,
		UserId:         userId,
	}
}
