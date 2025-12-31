package entity

import "time"

type CreateChatRoomEntity struct {
	CreateUserHash string
	RegDate        time.Time
	RoomKey        string
	RoomType       string
	Title          string
	SecretFlag     string
	Secret         string
	Description    string
	WorksCode      string
}

func MakeCreateChatRoomEntity(reqUserHash string, regDate time.Time, roomKey string, roomType string, title string, secretFlag string, secret string, des string, worksCode string) CreateChatRoomEntity {

	return CreateChatRoomEntity{
		CreateUserHash: reqUserHash,
		RegDate:        regDate,
		RoomKey:        roomKey,
		RoomType:       roomType,
		Title:          title,
		SecretFlag:     secretFlag,
		Secret:         secret,
		Description:    des,
		WorksCode:      worksCode,
	}

}
