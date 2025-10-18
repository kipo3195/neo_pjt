package entity

type RefreshTokenCheckEntity struct {
	UserId       string
	RefreshToken string
	Uuid         string
}

func MakeRefreshTokenCheckEntity(userId string, uuid string, rt string) RefreshTokenCheckEntity {
	return RefreshTokenCheckEntity{
		UserId:       userId,
		RefreshToken: rt,
		Uuid:         uuid,
	}
}
