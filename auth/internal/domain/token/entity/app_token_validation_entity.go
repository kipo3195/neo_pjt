package entity

type AppTokenValidationEntity struct {
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
}

func NewAppTokenValidationEntity(uuid string, token string) AppTokenValidationEntity {
	return AppTokenValidationEntity{
		Uuid:  uuid,
		Token: token,
	}
}
