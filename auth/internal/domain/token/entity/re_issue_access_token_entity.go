package entity

type ReIssueAccessTokenEntity struct {
	UserId string
	Uuid   string
}

func MakeReIssueAccessTokenEntity(userId string, uuid string) ReIssueAccessTokenEntity {
	return ReIssueAccessTokenEntity{
		UserId: userId,
		Uuid:   uuid,
	}
}
