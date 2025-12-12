package entity

type RegistServiceUserEntity struct {
	Org    string
	UserId []string
}

func MakeRegistServiceUserEntity(org string, userId []string) RegistServiceUserEntity {
	return RegistServiceUserEntity{
		Org:    org,
		UserId: userId,
	}
}
