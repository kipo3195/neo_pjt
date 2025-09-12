package entity

type UserAuthInfoEntity struct {
	Id       string `json:"id"`
	Salt     string `json:"salt"`
	AuthHash string `json:"authHash"`
	UserHash string `json:"userHash"`
}

func MakeUserAuthInfoEntity(id string, salt string, authHash string, userHash string) UserAuthInfoEntity {
	return UserAuthInfoEntity{
		Id:       id,
		Salt:     salt,
		UserHash: userHash,
		AuthHash: authHash,
	}
}
