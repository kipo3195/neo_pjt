package entity

type CreateChatRoomEntity struct {
	ChatRoomEntity       ChatRoomEntity               `json:"createChatRoom"`
	ChatRoomMemberEntity []CreateChatRoomMemberEntity `json:"createChatRoomMember"`
}
