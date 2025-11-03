package entity

type ProfileImgEntity struct {
	ProfileImg          *[]byte
	ProfileImgSize      int64
	ProfileImgName      string
	UserId              string
	UserHash            string
	ProfileImgSavedName string
	ProfileImgSavedPath string
	ProfileImgHash      string
}

func MakeProfileImgEntity(profileImg *[]byte, profileImgSize int64, profileImgName string, userId string, userHash string) ProfileImgEntity {
	return ProfileImgEntity{
		ProfileImg:     profileImg,
		ProfileImgSize: profileImgSize,
		ProfileImgName: profileImgName,
		UserId:         userId,
		UserHash:       userHash,
	}
}
