package chatRoom

import (
	"message/internal/delivery/dto/chatLine"
	"message/internal/delivery/dto/chatRoomConfig"
	"message/internal/delivery/dto/chatRoomFixed"
	"message/internal/delivery/dto/chatRoomTitle"
	"message/internal/delivery/dto/chatUnread"
)

type GetChatRoomListDto struct {
	ChatRoomDetail   ChatRoomDetail                `json:"roomDetail"`
	ChatRoomFixed    chatRoomFixed.ChatRoomFixed   `json:"fixed"`
	MyChatRoomTitle  chatRoomTitle.ChatRoomTitle   `json:"myChatRoomTitle"`
	MyChatRoomConfig chatRoomConfig.ChatRoomConfig `json:"myChatRoomConfig"`
	Member           []string                      `json:"member"`
	Owner            []string                      `json:"onwer"`
	Line             chatLine.ChatLineDto          `json:"line"`
	Unread           chatUnread.ChatUnreadDto      `json:"unread"`
}
