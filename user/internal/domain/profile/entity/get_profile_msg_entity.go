package entity

type GetProfileMsgEntity struct {
	UserHashs []string
}

func MakeGetProfileMsgEntity(userHashs []string) GetProfileMsgEntity {
	return GetProfileMsgEntity{
		UserHashs: userHashs,
	}
}
