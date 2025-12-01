package dto

import "notificator/internal/domain/socketSender/entity"

type SendChatDto struct {
	Type         string      `json:"type"`
	ChatSession  string      `json:"chatSession"`
	EventType    string      `json:"eventType"`
	ChatLineData ChatLineDto `json:"chatLineData"`
	ChatRoomData ChatRoomDto `json:"chatRoomData"`
}

func MakeSendChatDto(t string, eventType string, chatSession string, chatLineData entity.SendChatLineEntity, chatRoomData entity.SendChatRoomEntity) SendChatDto {

	chatLine := MakeChatLineDto(chatLineData.Cmd, chatLineData.Contents, chatLineData.LineKey, chatLineData.SendUserHash, chatLineData.SendDate)

	chatRoom := MakeChatRoomDto(chatRoomData.RoomType, chatRoomData.RoomKey, chatRoomData.SecretFlag)

	return SendChatDto{

		Type:         t,
		EventType:    eventType,
		ChatSession:  chatSession,
		ChatLineData: chatLine,
		ChatRoomData: chatRoom,
	}

}
