package entity

type UserAuthChallengeEntity struct {
	Id string
}

func MakeUserAuthChallengeEntity(id string) UserAuthChallengeEntity {
	return UserAuthChallengeEntity{
		Id: id,
	}
}
