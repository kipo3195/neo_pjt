package entity

type UserRegisterChallengeEntity struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
}

func MakeUserRegisterChallengeEntity(id string, salt string) UserRegisterChallengeEntity {

	return UserRegisterChallengeEntity{
		Id:   id,
		Salt: salt,
	}
}
