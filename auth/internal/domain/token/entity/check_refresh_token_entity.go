package entity

type CheckRefreshTokenEntity struct {
	UserId       string
	RefreshToken string
	Uuid         string
}

func MakeCheckRefreshTokenEntity(userId string, rt string, uuid string) CheckRefreshTokenEntity {
	return CheckRefreshTokenEntity{
		UserId:       userId,
		RefreshToken: rt,
		Uuid:         uuid,
	}
}
