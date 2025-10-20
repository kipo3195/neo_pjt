package entity

type UserInfoEntity struct {
	UserIds []string
}

func MakeUserInfoEntity(userIds []string) UserInfoEntity {
	return UserInfoEntity{
		UserIds: userIds,
	}
}
