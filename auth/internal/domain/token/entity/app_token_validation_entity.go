package entity

type AppTokenValidationEntity struct {
	Uuid     string `json:"uuid"`
	AppToken string `json:"appToken"`
	Token    string `json:"token"`
}

func NewAppTokenValidationEntity(uuid string, appToken string, token string) AppTokenValidationEntity {
	return AppTokenValidationEntity{
		Uuid:     uuid,
		AppToken: appToken,
		Token:    token,
	}
}
