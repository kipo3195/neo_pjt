package entity

type RegistServiceUserEntity struct {
	Org          string
	UserIdPrefix string
	UserId       []string
	Start        int
	End          int
}

func MakeRegistServiceUserEntity(org string, userId []string, userIdPrefix string, start int, end int) RegistServiceUserEntity {
	return RegistServiceUserEntity{
		Org:          org,
		UserId:       userId,
		UserIdPrefix: userIdPrefix,
		Start:        start,
		End:          end,
	}
}
