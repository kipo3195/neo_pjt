package entity

type ChatRoomEventEntity struct {
	ChatRoomEntity       ChatRoomEntity
	ChatRoomMemberEntity []ChatRoomMemberEntity
}

func MakeChatRoomEventEntity(chatRoom ChatRoomEntity, chatRoomMemberEntity []ChatRoomMemberEntity) ChatRoomEventEntity {

	return ChatRoomEventEntity{
		ChatRoomEntity:       chatRoom,
		ChatRoomMemberEntity: chatRoomMemberEntity,
	}

}
