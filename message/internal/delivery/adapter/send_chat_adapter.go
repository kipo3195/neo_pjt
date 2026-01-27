package adapter

import (
	"message/internal/application/usecase/input"
	"message/internal/delivery/dto/chatService"
)

func MakeSendChatInput(sendUserHash string, lineKey string, sendDate string, eventType string, chatSession string, chatRoom chatService.ChatRoomData, chatLine chatService.ChatLineData, chatFile []chatService.ChatFileData) input.SendChatInput {

	room := input.ChatRoomInput{
		RoomKey:    chatRoom.RoomKey,
		RoomType:   chatRoom.RoomType,
		SecretFlag: chatRoom.SecretFlag,
	}

	line := input.ChatLineInput{
		SendUserHash:  sendUserHash,
		LineKey:       lineKey,
		TargetLineKey: chatLine.TargetLineKey,
		Contents:      chatLine.Contents,
		Cmd:           chatLine.Cmd,
		SendDate:      sendDate,
	}

	file := make([]input.ChatFileInput, 0)

	for _, f := range chatFile {

		temp := input.ChatFileInput{
			FileId:   f.FileId,
			FileExt:  f.FileExt,
			FileName: f.FileName,
		}
		file = append(file, temp)
	}

	return input.SendChatInput{
		ChatRoom:    room,
		ChatLine:    line,
		ChatSession: chatSession,
		EventType:   eventType,
		ChatFile:    file,
	}
}
