package entity

type SubscribeChatEntity struct {
	UserHash string
}

func MakeSubscribeChatEntity(userHash string) SubscribeChatEntity {
	return SubscribeChatEntity{
		UserHash: userHash,
	}
}
