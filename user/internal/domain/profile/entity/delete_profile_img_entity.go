package entity

type DeleteProfileImgEntity struct {
	UserId string
}

func MakeDeleteProfileImgEntity(userId string) DeleteProfileImgEntity {
	return DeleteProfileImgEntity{
		UserId: userId,
	}
}
