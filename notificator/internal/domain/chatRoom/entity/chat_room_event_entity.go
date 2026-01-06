package entity

import "time"

type ChatRoomEventEntity struct {
	Type           string                 `json:"type"`
	EventType      string                 `json:"eventType"`
	ChatRoomData   ChatRoomDataEntity     `json:"chatRoomData"`
	ChatRoomMember []ChatRoomMemberEntity `json:"chatRoomMember"`
}

func MakeChatRoomEventEntity(t string, eventType string, reqUserHash string, regDate time.Time, roomKey string, roomType string, title string, secretFlag string, secret string, des string, worksCode string, chatRoomMemberEntity []ChatRoomMemberEntity) ChatRoomEventEntity {

	chatRoomDataEntity := ChatRoomDataEntity{
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

	return ChatRoomEventEntity{
		Type:           t,
		EventType:      eventType,
		ChatRoomData:   chatRoomDataEntity,
		ChatRoomMember: chatRoomMemberEntity,
	}

}
