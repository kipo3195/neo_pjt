package entity

type CreateChatRoomEntity struct {
	RoomKey              string
	RoomType             string
	CreateChatRoomMember []CreateChatRoomMemberEntity
}

func MakeCreateChatRoomEntity(roomKey string, roomType string, member []CreateChatRoomMemberEntity) CreateChatRoomEntity {
	return CreateChatRoomEntity{
		RoomKey:              roomKey,
		RoomType:             roomType,
		CreateChatRoomMember: member,
	}
}
