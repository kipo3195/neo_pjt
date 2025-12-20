package chatRoom

import (
	"message/internal/delivery/dto/chatRoomConfig"
	"message/internal/delivery/dto/chatRoomFixed"
	"message/internal/delivery/dto/chatRoomTitle"
)

type GetChatRoomListDto struct {
	ChatRoomDetail   ChatRoomDetail                `json:"roomDetail"`
	ChatRoomFixed    chatRoomFixed.ChatRoomFixed   `json:"fixed"`
	MyChatRoomTitle  chatRoomTitle.ChatRoomTitle   `json:"myChatRoomTitle"`
	MyChatRoomConfig chatRoomConfig.ChatRoomConfig `json:"myChatRoomConfig"`
	Member           []string                      `json:"member"`
}
