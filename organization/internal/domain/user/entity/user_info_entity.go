package entity

type UserInfoEntity struct {
	Userhashs []string
}

func MakeUserInfoEntity(userHashs []string) UserInfoEntity {
	return UserInfoEntity{
		Userhashs: userHashs,
	}
}
