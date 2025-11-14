package entity

type PutProfileMsgEntity struct {
	UserHash   string
	ProfileMsg string
}

func MakePutProfileMsgEntity(userHash string, profileMsg string) PutProfileMsgEntity {
	return PutProfileMsgEntity{
		UserHash:   userHash,
		ProfileMsg: profileMsg,
	}
}
