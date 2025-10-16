package entity

type CheckRefreshTokenEntity struct {
	UserId       string
	RefreshToken string
	Uuid         string
}

func MakeCheckRefreshTokenEntity(userId string, uuid string, rt string) CheckRefreshTokenEntity {
	return CheckRefreshTokenEntity{
		UserId:       userId,
		RefreshToken: rt,
		Uuid:         uuid,
	}
}
