package dto

import "notificator/internal/domain/socketSender/entity"

type ChatDto struct {
	Type         string      `json:"type"`
	ChatSession  string      `json:"chatSession"`
	EventType    string      `json:"eventType"`
	ChatLineData ChatLineDto `json:"chatLineData"`
	ChatRoomData ChatRoomDto `json:"chatRoomData"`
}

func MakeChatDto(t string, eventType string, chatSession string, chatLineData entity.ChatLineEntity, chatRoomData entity.ChatRoomEntity) ChatDto {

	chatLine := MakeChatLineDto(chatLineData.Cmd, chatLineData.Contents, chatLineData.LineKey, chatLineData.TargetLineKey, chatLineData.SendUserHash, chatLineData.SendDate)

	chatRoom := MakeChatRoomDto(chatRoomData.RoomType, chatRoomData.RoomKey, chatRoomData.SecretFlag)

	return ChatDto{

		Type:         t,
		EventType:    eventType,
		ChatSession:  chatSession,
		ChatLineData: chatLine,
		ChatRoomData: chatRoom,
	}

}
