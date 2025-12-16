package entity

type DeleteProfileImgEntity struct {
	UserHash string
}

func MakeDeleteProfileImgEntity(userHash string) DeleteProfileImgEntity {
	return DeleteProfileImgEntity{
		UserHash: userHash,
	}
}
