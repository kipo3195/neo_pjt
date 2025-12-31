package dto

import "notificator/internal/domain/socketSender/entity"

type SendChatResponse struct {
	Type         string          `json:"type"`
	ChatSession  string          `json:"chatSession"`
	EventType    string          `json:"eventType"`
	ChatLineData SendChatLineDto `json:"chatLineData"`
	ChatRoomData SendChatRoomDto `json:"chatRoomData"`
}

func MakeSendChatResponse(t string, eventType string, chatSession string, chatLineData entity.ChatLineEntity, chatRoomData entity.ChatRoomEntity) SendChatResponse {

	chatLine := MakeChatLineDto(chatLineData.Cmd, chatLineData.Contents, chatLineData.LineKey, chatLineData.TargetLineKey, chatLineData.SendUserHash, chatLineData.SendDate)

	chatRoom := MakeChatRoomDto(chatRoomData.RoomType, chatRoomData.RoomKey, chatRoomData.SecretFlag)

	return SendChatResponse{
		Type:         t,
		EventType:    eventType,
		ChatSession:  chatSession,
		ChatLineData: chatLine,
		ChatRoomData: chatRoom,
	}

}
