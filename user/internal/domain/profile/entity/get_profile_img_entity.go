package entity

type GetProfileImgEntity struct {
	UserHash string
}

func MakeGetProfileImgEntity(userHash string) GetProfileImgEntity {
	return GetProfileImgEntity{
		UserHash: userHash,
	}
}
