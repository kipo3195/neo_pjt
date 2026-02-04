package chatRoom

import (
	"message/internal/adapter/http/dto/chatLine"
	"message/internal/adapter/http/dto/chatRoomConfig"
	"message/internal/adapter/http/dto/chatRoomFixed"
	"message/internal/adapter/http/dto/chatRoomTitle"
	"message/internal/adapter/http/dto/chatUnread"
)

type GetChatRoomDetailDto struct {
	ChatRoomDetail   ChatRoomDetail                `json:"roomDetail"`
	ChatRoomFixed    chatRoomFixed.ChatRoomFixed   `json:"fixed"`
	MyChatRoomTitle  chatRoomTitle.ChatRoomTitle   `json:"myChatRoomTitle"`
	MyChatRoomConfig chatRoomConfig.ChatRoomConfig `json:"myChatRoomConfig"`
	Member           []string                      `json:"member"`
	Owner            []string                      `json:"owner"`
	Line             chatLine.ChatLineDto          `json:"line"`
	Unread           chatUnread.ChatUnreadDto      `json:"unread"`
}
