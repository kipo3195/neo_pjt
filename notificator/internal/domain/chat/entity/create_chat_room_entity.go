package entity

type CreateChatRoomEntity struct {
	RoomKey              string
	RoomType             string
	CreateChatRoomMember []ChatRoomMemberEntity
}

func MakeCreateChatRoomEntity(roomKey string, roomType string, member []ChatRoomMemberEntity) CreateChatRoomEntity {
	return CreateChatRoomEntity{
		RoomKey:              roomKey,
		RoomType:             roomType,
		CreateChatRoomMember: member,
	}
}
