package entity

type GetUserDetailInfoEntity struct {
	UserHashs []string
}

func MakeGetUserDetailInfoEntity(userHashs []string) GetUserDetailInfoEntity {
	return GetUserDetailInfoEntity{
		UserHashs: userHashs,
	}
}
