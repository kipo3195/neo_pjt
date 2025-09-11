package entity

type UserAuthEntity struct {
	Id       string `json:"id"`
	Salt     string `json:"salt"`
	AuthHash string `json:"authHash"`
	UserHash string `json:"userHash"`
}

func MakeUserAuthEntity(id string, salt string, authHash string, userHash string) UserAuthEntity {
	return UserAuthEntity{
		Id:       id,
		Salt:     salt,
		UserHash: userHash,
		AuthHash: authHash,
	}
}
