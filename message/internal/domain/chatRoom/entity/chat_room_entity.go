package entity

import "time"

type ChatRoomEntity struct {
	CreateUserHash string    `json:"createUserHash"`
	RegDate        time.Time `json:"regDate"`
	RoomKey        string    `json:"roomKey"`
	RoomType       string    `json:"roomType"`
	Title          string    `json:"title"`
	SecretFlag     string    `json:"secretFlag"`
	Secret         string    `json:"secret"`
	Description    string    `json:"description"`
	WorksCode      string    `json:"worksCode"`
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
