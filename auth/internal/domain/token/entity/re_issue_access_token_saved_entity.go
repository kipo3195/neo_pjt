package entity

type ReIssueAccessTokenSavedEntity struct {
	UserId string
	Uuid   string
	Rt     string
	At     string
}

func MakeReIssueAccessTokenSavedEntity(userId string, uuid string, rt string, at string) ReIssueAccessTokenSavedEntity {
	return ReIssueAccessTokenSavedEntity{
		UserId: userId,
		Uuid:   uuid,
		Rt:     rt,
		At:     at,
	}
}
