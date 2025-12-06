package entity

import "time"

type ChatRoomEntity struct {
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

// 생성과 조회 entity는 동일하나 이벤트에 따른 구분 처리를 위함
func MakeCreateChatRoomEntity(createUserHash string, regDate time.Time, roomKey string, roomType string, title string, description string, secretFlag string, secret string, worksCode string) ChatRoomEntity {

	return ChatRoomEntity{
		CreateUserHash: createUserHash,
		RegDate:        regDate,
		RoomKey:        roomKey,
		RoomType:       roomType,
		Title:          title,
		Description:    description,
		SecretFlag:     secretFlag,
		Secret:         secret,
		WorksCode:      worksCode,
	}
}
