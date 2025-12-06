package adapter

import (
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"message/internal/delivery/dto/chatRoom"
)

func MakeCreateChatRoomInput(reqUser string, roomKey string, roomType string, title string, secretFlag string, secret string, description string, worksCode string, member []chatRoom.ChatRoomMember) (input.CreateChatRoomInput, error) {

	memberInput := make([]input.CreateChatMemberInput, 0)

	var regFlag = false
	for _, m := range member {

		temp := input.CreateChatMemberInput{
			MemberHash:      m.MemberHash,
			MemberWorksCode: m.MemberWorksCode,
		}

		memberInput = append(memberInput, temp)
		if reqUser == temp.MemberHash {
			regFlag = true
		}
	}

	// 생성자가 참여자에 없으면 error
	if !regFlag {
		return input.CreateChatRoomInput{}, consts.ErrInvalidChatRoomMember
	}

	return input.CreateChatRoomInput{
		CreateUserHash: reqUser,
		RoomKey:        roomKey,
		RoomType:       roomType,
		Title:          title,
		SecretFlag:     secretFlag,
		Secret:         secret,
		Description:    description,
		WorksCode:      worksCode,
		Member:         memberInput,
	}, nil

}
