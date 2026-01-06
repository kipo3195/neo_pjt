package input

type ChatRoomEventInput struct {
	Type                     string                     `json:"type"`
	EventType                string                     `json:"eventType"`
	ChatRoomEventDataInput   ChatRoomEventDataInput     `json:"chatRoomData"`
	ChatRoomEventMemberInput []ChatRoomEventMemberInput `json:"chatRoomMember"`
}
