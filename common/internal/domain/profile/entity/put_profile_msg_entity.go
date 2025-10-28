package entity

type PutProfileMsgEntity struct {
	UserId string
	Msg    string
}

func MakePutProfileMsgEntity(userId string, msg string) PutProfileMsgEntity {
	return PutProfileMsgEntity{
		UserId: userId,
		Msg:    msg,
	}
}
