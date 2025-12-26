package input

type CreateChatRoomMessageInput struct {
	CreateChatRoomInput       CreateChatRoomInput         `json:"createChatRoom"`
	CreateChatRoomMemberInput []CreateChatRoomMemberInput `json:"createChatRoomMember"`
}
