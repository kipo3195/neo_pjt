package entity

type GetProfileImgEntity struct {
	UserId string
}

func MakeGetProfileImgEntity(userId string) GetProfileImgEntity {
	return GetProfileImgEntity{
		UserId: userId,
	}
}
