package dto

import "time"

type CreateChatRoomResponse struct {
	Type              string            `json:"type"`
	EventType         string            `json:"eventType"`
	CreateChatRoomDto CreateChatRoomDto `json:"chatRoomData"`
}

func MakeCreateChatRoomResponse(reqUserHash string, regDate time.Time, roomKey string, roomType string, title string, secretFlag string, secret string, des string, worksCode string) CreateChatRoomResponse {

	createChatRoomDto := CreateChatRoomDto{
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

	return CreateChatRoomResponse{
		Type:              "chatRoom",
		EventType:         "C",
		CreateChatRoomDto: createChatRoomDto,
	}

}
