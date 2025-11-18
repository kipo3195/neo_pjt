package entity

type SubscribeNoteEntity struct {
	UserHash string
}

func MakeSubscribeNoteEntity(userHash string) SubscribeNoteEntity {
	return SubscribeNoteEntity{
		UserHash: userHash,
	}
}
