package output

type GetProfileImgOutput struct {
	ProfileImg     []byte
	ProfileImgName string
}

func MakeGetProfileImgOutput(profileImg []byte, profileImgName string) GetProfileImgOutput {
	return GetProfileImgOutput{
		ProfileImg:     profileImg,
		ProfileImgName: profileImgName,
	}

}
